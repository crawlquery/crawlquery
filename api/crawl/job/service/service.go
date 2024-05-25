package service

import (
	"crawlquery/api/domain"
	"crawlquery/node/dto"
	"crawlquery/pkg/util"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo               domain.CrawlJobRepository
	shardService       domain.ShardService
	nodeService        domain.NodeService
	restrictionService domain.CrawlRestrictionService
	pageService        domain.PageService
	logger             *zap.SugaredLogger
}

func NewService(
	repo domain.CrawlJobRepository,
	shardService domain.ShardService,
	nodeService domain.NodeService,
	restrictionService domain.CrawlRestrictionService,
	pageService domain.PageService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		repo:               repo,
		shardService:       shardService,
		nodeService:        nodeService,
		restrictionService: restrictionService,
		pageService:        pageService,
		logger:             logger,
	}
}

var trackingParams = []string{
	"utm_source", "utm_medium", "utm_platform", "utm_campaign", "utm_term", "utm_content", "gclid", "fbclid",
}

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	queryParams := parsedURL.Query()
	for _, param := range trackingParams {
		queryParams.Del(param)
	}
	parsedURL.RawQuery = queryParams.Encode()

	parsedURL.Host = strings.ToLower(parsedURL.Host)
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)

	return parsedURL.String(), nil
}

func (cs *Service) Create(url string) (*domain.CrawlJob, error) {

	normalizedURL, err := normalizeURL(url)

	if err != nil {
		cs.logger.Errorw("Crawl.Service.AddJob: error normalizing URL", "error", err)
		return nil, err
	}

	job := &domain.CrawlJob{
		ID:        util.UUID(),
		URL:       normalizedURL,
		PageID:    util.PageID(url),
		CreatedAt: time.Now(),
	}

	if err := job.Validate(); err != nil {
		return nil, err
	}

	// Save the job in the repository
	if err := cs.repo.Create(job); err != nil {
		cs.logger.Errorw("Crawl.Service.AddJob: error creating job", "error", err)
		return nil, err
	}
	return job, nil
}

func (cs *Service) pushBack(job *domain.CrawlJob, until time.Time, reason string) error {
	cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error processing job", "error", reason, "job", job)
	job.BackoffUntil = sql.NullTime{
		Time:  until,
		Valid: true,
	}
	job.FailedReason = sql.NullString{
		String: reason,
		Valid:  true,
	}

	if err := cs.repo.Update(job); err != nil {
		cs.logger.Errorw("Crawl.Service.pushBack: error updating job", "error", reason)
		return err
	}

	return nil
}

func (cs *Service) ProcessCrawlJobs() {
	for {
		// Get the first job
		job, err := cs.repo.FirstProcessable()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// If there are no jobs, wait and try again
		if job == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		pasedURL, err := url.Parse(job.URL)

		if err != nil {
			cs.pushBack(job, time.Now().Add(time.Hour), err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		// Check if the domain is restricted
		restricted, until := cs.restrictionService.GetRestriction(pasedURL.Hostname())

		if restricted {
			cs.pushBack(job, *until, fmt.Sprintf("domain is restricted until %v", until))
			continue
		}

		shardID, err := cs.shardService.GetURLShardID(job.URL)

		if err != nil {
			cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error getting shardID", "error", err)
			continue
		}

		// Process the job
		crawledPage, err := cs.processJob(job, shardID)

		if err != nil {
			cs.pushBack(job, time.Now().Add(time.Hour), err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		job.LastCrawledAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		if err := cs.repo.Update(job); err != nil {
			cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error updating job", "error", err)
			time.Sleep(5 * time.Second)
		}

		_, err = cs.pageService.Create(job.PageID, shardID, crawledPage.Hash)

		if err != nil {
			cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error creating page", "error", err)
			continue
		}

		cs.restrictionService.Restrict(pasedURL.Hostname())
	}
}

func (cs *Service) processJob(job *domain.CrawlJob, shardID uint) (*dto.Page, error) {
	// Process the job
	cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: processing job", "job", job)

	nodes, err := cs.nodeService.ListByShardID(shardID)

	if err != nil {
		cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error getting nodes", "error", err)
		return nil, err
	}

	nodes = cs.nodeService.Randomize(nodes)

	cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: nodes", "nodes len", len(nodes))

	if len(nodes) == 0 {
		cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: no nodes available", "nodes", nodes)
		return nil, errors.New("no nodes available")
	}

	cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: processing node", "node", nodes[0])

	// Send the job to the node
	crawledPage, err := cs.nodeService.SendCrawlJob(nodes[0], job)

	if err != nil {
		cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error sending job to node", "error", err)
		return nil, err
	}

	cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: job sent to node", "node", nodes[0])
	return crawledPage, nil
}

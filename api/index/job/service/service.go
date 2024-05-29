package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	indexJobRepo domain.IndexJobRepository
	pageService  domain.PageService
	nodeService  domain.NodeService
	logger       *zap.SugaredLogger
}

func NewService(
	indexJobRepo domain.IndexJobRepository,
	pageService domain.PageService,
	nodeService domain.NodeService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		indexJobRepo: indexJobRepo,
		pageService:  pageService,
		nodeService:  nodeService,
		logger:       logger,
	}
}

func (s *Service) Create(pageID string) (*domain.IndexJob, error) {

	if _, err := s.indexJobRepo.GetByPageID(pageID); err == nil {
		return nil, domain.ErrIndexJobAlreadyExists
	}

	job := &domain.IndexJob{
		ID:        util.UUIDString(),
		PageID:    pageID,
		CreatedAt: time.Now(),
	}

	return s.indexJobRepo.Create(job)
}

func (s *Service) Get(id string) (*domain.IndexJob, error) {
	return s.indexJobRepo.Get(id)
}

func (s *Service) Next() (*domain.IndexJob, error) {
	return s.indexJobRepo.Next()
}

func (s *Service) Update(job *domain.IndexJob) error {
	return s.indexJobRepo.Update(job)
}

func (s *Service) processJob(job *domain.IndexJob, shardID uint) error {
	// Process the job
	s.logger.Infow("Crawl.Service.ProcessCrawlJobs: processing job", "job", job)

	nodes, err := s.nodeService.ListByShardID(shardID)

	if err != nil {
		s.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error getting nodes", "error", err)
		return err
	}

	nodes = s.nodeService.Randomize(nodes)

	s.logger.Infow("Crawl.Service.ProcessCrawlJobs: nodes", "nodes len", len(nodes))

	if len(nodes) == 0 {
		s.logger.Errorw("Crawl.Service.ProcessCrawlJobs: no nodes available", "nodes", nodes)
		return errors.New("no nodes available")
	}

	s.logger.Infow("Crawl.Service.ProcessCrawlJobs: processing node", "node", nodes[0])

	// Send the job to the node
	err = s.nodeService.SendIndexJob(nodes[0], job)

	if err != nil {
		s.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error sending job to node", "error", err)
		return err
	}

	s.logger.Infow("Crawl.Service.ProcessCrawlJobs: job complete", "node", nodes[0])
	return nil
}

func (s *Service) ProcessIndexJobs() {
	for {
		job, err := s.Next()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if job == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		s.logger.Infow("Processing index job", "job", job)

		page, err := s.pageService.Get(job.PageID)

		if err != nil {
			job.BackoffUntil = sql.NullTime{
				Time:  time.Now().Add(24 * time.Hour),
				Valid: true,
			}

			job.FailedReason = sql.NullString{
				String: "page not found",
				Valid:  true,
			}

			if err := s.Update(job); err != nil {
				s.logger.Errorw("Error updating index job", "error", err)
			}
		}

		err = s.processJob(job, page.ShardID)

		if err != nil {
			job.BackoffUntil = sql.NullTime{
				Time:  time.Now().Add(24 * time.Hour),
				Valid: true,
			}

			job.FailedReason = sql.NullString{
				String: err.Error(),
				Valid:  true,
			}

			if err := s.Update(job); err != nil {
				s.logger.Errorw("Error updating index job", "error", err)
			}

			continue
		}

		job.LastIndexedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		if err := s.Update(job); err != nil {
			s.logger.Errorw("Error updating index job", "error", err)
		}

		s.logger.Infow("Index job processed", "job", job)
	}
}

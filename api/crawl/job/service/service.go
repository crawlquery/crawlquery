package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"errors"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo         domain.CrawlJobRepository
	shardService domain.ShardService
	nodeService  domain.NodeService
	logger       *zap.SugaredLogger
}

func NewService(
	repo domain.CrawlJobRepository,
	shardService domain.ShardService,
	nodeService domain.NodeService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		repo:         repo,
		shardService: shardService,
		nodeService:  nodeService,
		logger:       logger,
	}
}

func (cs *Service) Create(url string) (*domain.CrawlJob, error) {
	job := &domain.CrawlJob{
		ID:        util.UUID(),
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := job.Validate(); err != nil {
		return nil, err
	}

	// Save the job in the repository
	if err := cs.repo.Create(job); err != nil {
		cs.logger.Errorw("Crawl.Service.AddJob: error creating job", "error", err)
		return nil, domain.ErrInternalError
	}
	return job, nil
}

func (cs *Service) ProcessCrawlJobs() {
	for {
		// Get the first job
		job, err := cs.repo.First()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// If there are no jobs, wait and try again
		if job == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Process the job
		err = cs.processJob(job)

		if err != nil {
			cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error processing job", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Delete the job
		if err := cs.repo.Delete(job.ID); err != nil {
			cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error deleting job", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

func (cs *Service) processJob(job *domain.CrawlJob) error {
	// Process the job
	cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: processing job", "job", job)

	shardID, err := cs.shardService.GetURLShardID(job.URL)

	if err != nil {
		cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error getting shardID", "error", err)
		return err
	}

	cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: shardID", "shardID", shardID)

	nodes, err := cs.nodeService.ListByShardID(shardID)

	if err != nil {
		cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error getting nodes", "error", err)
		return err
	}

	cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: nodes", "nodes", nodes)

	for _, node := range nodes {
		cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: processing node", "node", node)

		// Send the job to the node
		err := cs.nodeService.SendCrawlJob(node, job)

		if err != nil {
			cs.logger.Errorw("Crawl.Service.ProcessCrawlJobs: error sending job to node", "error", err)
			continue
		}

		cs.logger.Infow("Crawl.Service.ProcessCrawlJobs: job sent to node", "node", node)
		cs.repo.Delete(job.ID)

		return nil
	}

	return errors.New("no nodes available")
}

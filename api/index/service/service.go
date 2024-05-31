package service

import (
	"context"
	"crawlquery/api/domain"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	eventService domain.EventService
	indexJobRepo domain.IndexJobRepository
	indexLogRepo domain.IndexLogRepository
	pageService  domain.PageService
	nodeService  domain.NodeService
	logger       *zap.SugaredLogger
	workers      int
	maxQueueSize int
}

type Option func(*Service)

func WithEventService(eventService domain.EventService) Option {
	return func(s *Service) {
		s.eventService = eventService
	}
}

func WithEventListeners() Option {
	return func(s *Service) {
		s.registerEventListeners()
	}
}

func WithIndexJobRepo(indexJobRepo domain.IndexJobRepository) Option {
	return func(s *Service) {
		s.indexJobRepo = indexJobRepo
	}
}

func WithIndexLogRepo(indexLogRepo domain.IndexLogRepository) Option {
	return func(s *Service) {
		s.indexLogRepo = indexLogRepo
	}
}

func WithPageService(pageService domain.PageService) Option {
	return func(s *Service) {
		s.pageService = pageService
	}
}

func WithNodeService(nodeService domain.NodeService) Option {
	return func(s *Service) {
		s.nodeService = nodeService
	}
}

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithWorkers(workers int) Option {
	return func(s *Service) {
		s.workers = workers
	}
}

func WithMaxQueueSize(maxQueueSize int) Option {
	return func(s *Service) {
		s.maxQueueSize = maxQueueSize
	}
}

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) registerEventListeners() {
	if s.eventService == nil {
		s.logger.Fatal("EventService is required")
	}

	s.eventService.Subscribe(domain.CrawlCompletedKey, s.handleCrawlCompleted)
}

func (s *Service) handleCrawlCompleted(event domain.Event) {
	crawlCompleted := event.(*domain.CrawlCompleted)

	err := s.CreateJob(crawlCompleted.PageID, crawlCompleted.ShardID)

	if err != nil {
		s.logger.Errorf("Error creating index job: %v", err)
	}

}

func (s *Service) createlogEntry(pageID domain.PageID, status domain.IndexStatus) error {
	log := &domain.IndexLog{
		PageID:    pageID,
		Status:    status,
		CreatedAt: time.Now(),
	}

	err := s.indexLogRepo.Save(log)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateJob(pageID domain.PageID, shardID domain.ShardID) error {
	if _, err := s.indexJobRepo.Get(pageID); err == nil {
		return domain.ErrIndexJobAlreadyExists
	}

	job := &domain.IndexJob{
		PageID:    pageID,
		Status:    domain.IndexStatusPending,
		CreatedAt: time.Now(),
	}

	err := s.indexJobRepo.Save(job)
	if err != nil {
		return err
	}

	return s.createlogEntry(pageID, domain.IndexStatusPending)
}

func (s *Service) RunIndexProcess(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := s.processJobsWithWorkers(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (s *Service) processJobsWithWorkers(ctx context.Context) error {
	jobs := make(chan *domain.IndexJob)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.processJob(ctx, jobs, &wg)
	}

	jobsToProcess, err := s.indexJobRepo.ListByStatus(10000, domain.IndexStatusPending)
	if err != nil {
		return err
	}

	for _, job := range jobsToProcess {
		job.Status = domain.IndexStatusInProgress
		err = s.indexJobRepo.Save(job)
		if err != nil {
			s.logger.Errorf("Error saving job: %v", err)
			continue
		}
		jobs <- job
	}

	close(jobs)
	wg.Wait()
	return nil
}

func (s *Service) processJob(ctx context.Context, jobs <-chan *domain.IndexJob, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes, err := s.nodeService.RandomizedListGroupByShard()
	if err != nil {
		s.logger.Errorf("processJob failed to get nodes: %v", err)
		return
	}

	for job := range jobs {
		jobCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		var maxAttempts = 3
		var attempts = 0
		for _, node := range nodes[job.ShardID] {
			attempts++
			err := s.indexPage(jobCtx, job, node)

			if err == nil {
				job.Status = domain.IndexStatusCompleted
				s.indexJobRepo.Save(job)
				s.createlogEntry(job.PageID, domain.IndexStatusCompleted)
				break
			}

			if attempts >= maxAttempts {
				s.logger.Errorw("Failed to process job after max attempts", "job", job, "error", err)
				job.Status = domain.IndexStatusFailed
				s.indexJobRepo.Save(job)
				s.createlogEntry(job.PageID, domain.IndexStatusFailed)
				break
			}
		}
	}
}

func (s *Service) indexPage(ctx context.Context, job *domain.IndexJob, node *domain.Node) error {
	s.logger.Infof("Indexing page %s on node %s", job.PageID, node.ID)
	return s.nodeService.SendIndexJob(ctx, node, job)
}
package service

import (
	"context"
	"crawlquery/api/domain"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	eventService         domain.EventService
	crawlJobRepo         domain.CrawlJobRepository
	crawlLogRepo         domain.CrawlLogRepository
	crawlThrottleService domain.CrawlThrottleService

	nodeService domain.NodeService

	logger *zap.SugaredLogger

	workers      int
	maxQueueSize int
}

type Option func(*Service)

func WithEventService(eventService domain.EventService) func(*Service) {
	return func(s *Service) {
		s.eventService = eventService
	}
}

func WithCrawlThrottleService(crawlThrottleService domain.CrawlThrottleService) func(*Service) {
	return func(s *Service) {
		s.crawlThrottleService = crawlThrottleService
	}
}

func WithCrawlJobRepo(crawlJobRepo domain.CrawlJobRepository) func(*Service) {
	return func(s *Service) {
		s.crawlJobRepo = crawlJobRepo
	}
}

func WithCrawlLogRepo(crawlLogRepo domain.CrawlLogRepository) func(*Service) {
	return func(s *Service) {
		s.crawlLogRepo = crawlLogRepo
	}
}

func WithLogger(logger *zap.SugaredLogger) func(*Service) {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithNodeService(nodeService domain.NodeService) func(*Service) {
	return func(s *Service) {
		s.nodeService = nodeService
	}
}

func WithWorkers(workers int) func(*Service) {
	return func(s *Service) {
		s.workers = workers
	}
}

func WithMaxQueueSize(maxQueueSize int) func(*Service) {
	return func(s *Service) {
		s.maxQueueSize = maxQueueSize
	}
}

func WithEventListeners() func(*Service) {
	return func(s *Service) {
		s.registerEventListeners()
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
	s.eventService.Subscribe(domain.PageCreatedKey, s.handlePageCreated)
}

func (s *Service) handlePageCreated(e domain.Event) {
	page := e.(*domain.PageCreated).Page

	err := s.CreateJob(page)

	if err != nil {
		s.logger.Errorw("Error creating job", "error", err)
	}

}

func (s *Service) CreateJob(page *domain.Page) error {

	if _, err := s.crawlJobRepo.Get(page.ID); err == nil {
		return nil
	}

	cj := &domain.CrawlJob{
		PageID: page.ID,
		URL:    page.URL,
		Status: domain.CrawlStatusPending,
	}

	err := s.crawlJobRepo.Save(cj)
	if err != nil {
		return err
	}

	cl := &domain.CrawlLog{
		PageID:    cj.PageID,
		Status:    domain.CrawlStatusPending,
		CreatedAt: time.Now(),
	}

	err = s.crawlLogRepo.Save(cl)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) jobsToProcess() ([]*domain.CrawlJob, error) {
	jobs, err := s.crawlJobRepo.ListByStatus(10000, domain.CrawlStatusPending)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *Service) RunCrawlProcess(ctx context.Context) error {
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

// processJobsWithWorkers processes jobs in the queue using a worker pool
func (s *Service) processJobsWithWorkers(ctx context.Context) error {
	jobs := make(chan *domain.CrawlJob)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.processJob(ctx, jobs, &wg)
	}

	jobsToProcess, err := s.jobsToProcess()

	if err != nil {
		return err
	}

	for _, job := range jobsToProcess {

		canCrawl, err := s.crawlThrottleService.CheckAndThrottle(job.URL)

		if err != nil {
			s.logger.Errorw("Throttle returned an error", "error", err)
			err = s.updateJob(job, domain.CrawlStatusFailed, err)
			if err != nil {
				return err
			}
			continue
		}

		if !canCrawl {
			s.logger.Infow("Throttling", "url", job.URL)
			err = s.updateJob(job, domain.CrawlStatusPending, err)
			if err != nil {
				return err
			}
			continue
		}

		err = s.updateJob(job, domain.CrawlStatusInProgress, err)
		if err != nil {
			return err
		}
		jobs <- job
	}

	close(jobs)
	wg.Wait()
	return nil
}

func (s *Service) processJob(ctx context.Context, jobs <-chan *domain.CrawlJob, wg *sync.WaitGroup) {
	nodes, err := s.nodeService.RandomizedListGroupByShard()

	if err != nil {
		s.logger.Error("processJob failed to get nodes: %v", err)
		return
	}

	defer wg.Done()
	for job := range jobs {
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		var maxAttempts = 3
		var attempts = 0
		for _, node := range nodes[job.ShardID] {
			attempts++
			err := s.ProcessQueueItem(ctx, job, node)

			if err == nil {
				break
			}

			if attempts >= maxAttempts {
				s.logger.Errorw("Failed to process job after max attempts", "job", job, "error", err)
				break
			}
		}
	}
}

func (s *Service) updateJob(job *domain.CrawlJob, status domain.CrawlStatus, withErr error) error {
	job.Status = status
	job.UpdatedAt = time.Now()
	err := s.crawlJobRepo.Save(job)
	if err != nil {
		return err
	}

	if withErr == nil {
		return s.updateLog(job, status, "")
	}

	return s.updateLog(job, status, withErr.Error())
}

func (s *Service) updateLog(job *domain.CrawlJob, status domain.CrawlStatus, info string) error {
	cl := &domain.CrawlLog{
		PageID:    job.PageID,
		Status:    status,
		Info:      info,
		CreatedAt: time.Now(),
	}
	return s.crawlLogRepo.Save(cl)
}

func (s *Service) ProcessQueueItem(ctx context.Context, job *domain.CrawlJob, assignedNode *domain.Node) error {
	cl := &domain.CrawlLog{
		PageID:    job.PageID,
		Status:    domain.CrawlStatusInProgress,
		CreatedAt: time.Now(),
	}

	err := s.crawlLogRepo.Save(cl)
	if err != nil {
		s.logger.Errorw("Error saving crawl log", "error", err)
		return err
	}

	res, err := s.nodeService.SendCrawlJob(
		ctx,
		assignedNode,
		job,
	)

	if err != nil {
		if err := s.updateJob(job, domain.CrawlStatusFailed, err); err != nil {
			return err
		}

		return err
	}

	s.updateJob(job, domain.CrawlStatusCompleted, nil)

	var links []domain.URL
	for _, link := range res.Links {
		links = append(links, domain.URL(link))
	}

	s.eventService.Publish(&domain.CrawlCompleted{
		PageID:      job.PageID,
		ShardID:     job.ShardID,
		ContentHash: domain.ContentHash(res.ContentHash),
		Links:       links,
	})

	return nil
}

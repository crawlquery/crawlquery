package service

import (
	"context"
	"crawlquery/api/domain"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	eventService domain.EventService
	crawlJobRepo domain.CrawlJobRepository
	crawlLogRepo domain.CrawlLogRepository
	crawlQueue   domain.CrawlQueue
	linkService  domain.LinkService

	nodeService domain.NodeService
	nodesCache  map[domain.ShardID][]*domain.Node

	logger *zap.SugaredLogger

	workers int
}

type Option func(*Service)

func WithEventService(eventService domain.EventService) func(*Service) {
	return func(s *Service) {
		s.eventService = eventService
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

func WithCrawlQueue(crawlQueue domain.CrawlQueue) func(*Service) {
	return func(s *Service) {
		s.crawlQueue = crawlQueue
	}
}

func WithLinkService(linkService domain.LinkService) func(*Service) {
	return func(s *Service) {
		s.linkService = linkService
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

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
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

func (s *Service) FillQueue() error {
	jobs, err := s.crawlJobRepo.ListByStatus(10000, domain.CrawlStatusPending)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		err := s.crawlQueue.Push(job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) cacheNodes() error {
	var err error

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Crawl(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

			err := s.FillQueue()
			if err != nil {
				return err
			}

			err = s.ProcessQueue()
			if err != nil {
				return err
			}
		}
	}
}

// ProcessQueue processes jobs in the queue using a worker pool
func (s *Service) ProcessQueue() error {
	jobs := make(chan *domain.CrawlJob)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.worker(jobs, &wg)
	}

	// Fetch jobs from the queue and send to workers
	for {
		job, err := s.crawlQueue.Pop()
		if err != nil {
			if errors.Is(err, domain.ErrCrawlQueueEmpty) {
				break // Exit loop if queue is empty
			}
			return err
		}
		jobs <- job
	}
	close(jobs)
	wg.Wait()
	return nil
}

func (s *Service) worker(jobs <-chan *domain.CrawlJob, wg *sync.WaitGroup) {
	nodes, err := s.nodeService.RandomizedListGroupByShard()

	if err != nil {
		s.logger.Error("worker failed to get nodes: %v", err)
		return
	}

	defer wg.Done()
	for job := range jobs {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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

func (s *Service) updateJob(job *domain.CrawlJob, status domain.CrawlStatus) error {
	job.Status = status
	job.UpdatedAt = time.Now()
	return s.crawlJobRepo.Save(job)
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
		if err := s.updateJob(job, domain.CrawlStatusFailed); err != nil {
			return err
		}
		s.updateLog(job, domain.CrawlStatusFailed, err.Error())
		return err
	}

	if err := s.updateJob(job, domain.CrawlStatusCompleted); err != nil {
		return err
	}

	var links []domain.URL
	for _, link := range res.Links {
		links = append(links, domain.URL(link))
	}

	s.eventService.Publish(&domain.CrawlCompleted{
		PageID: job.PageID,
		Links:  links,
	})

	s.updateLog(job, domain.CrawlStatusCompleted, "")

	return nil
}

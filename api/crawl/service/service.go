package service

import (
	"context"
	"crawlquery/api/domain"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Service struct {
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

func (s *Service) CacheNodes() error {
	var err error
	s.nodesCache, err = s.nodeService.RandomizedListGroupByShard()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetRandomNodeFromCache(shardID domain.ShardID) *domain.Node {
	nodes, ok := s.nodesCache[shardID]
	if !ok {
		return nil
	}

	if len(nodes) == 0 {
		return nil
	}

	var list []*domain.Node
	for _, node := range nodes {
		list = append(list, node)
	}

	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})

	return list[0]
}

func (s *Service) Crawl(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

			err := s.CacheNodes()
			if err != nil {
				return err
			}

			err = s.FillQueue()
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
	defer wg.Done()
	for job := range jobs {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := s.ProcessQueueItem(ctx, job)
		if err != nil {
			log.Printf("Failed to process job %s: %v", job.PageID, err)
		}
	}
}

func (s *Service) ProcessQueueItem(ctx context.Context, job *domain.CrawlJob) error {
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

	// Simulate processing time
	select {
	case <-ctx.Done():

		cl.Status = domain.CrawlStatusFailed
		cl.Info = "Crawl timed out"
		err := s.crawlLogRepo.Save(cl)

		if err != nil {
			s.logger.Errorw("Error saving crawl log", "error", err)
			return err
		}

		return ctx.Err()
	default:
		randomNode := s.GetRandomNodeFromCache(job.ShardID)
		if randomNode == nil {
			cl.Status = domain.CrawlStatusFailed
			cl.Info = "No nodes available"
			err := s.crawlLogRepo.Save(cl)
			if err != nil {
				s.logger.Errorw("Error saving crawl log", "error", err)
				return err
			}
			return nil
		}

		res, err := s.nodeService.SendCrawlJob(randomNode, job)

		if err != nil {
			cl.Status = domain.CrawlStatusFailed
			cl.Info = err.Error()
			err := s.crawlLogRepo.Save(cl)
			if err != nil {
				s.logger.Errorw("Error saving crawl log", "error", err)
				return err
			}
			return nil
		}

		for _, link := range res.Links {
			_, err := s.linkService.Create(job.PageID, domain.URL(link))
			if err != nil {
				s.logger.Errorw("Error creating crawl job", "error", err)
				return err
			}
		}

		job.Status = domain.CrawlStatusCompleted
		err = s.crawlJobRepo.Save(job)
		if err != nil {
			s.logger.Errorw("Error saving crawl job", "error", err)
			return err
		}
	}

	cl.Status = domain.CrawlStatusCompleted
	cl.CreatedAt = time.Now()
	return s.crawlLogRepo.Save(cl)
}

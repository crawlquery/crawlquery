package service

import (
	"crawlquery/api/domain"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	indexJobRepo domain.IndexJobRepository
	indexLogRepo domain.IndexLogRepository
	pageService  domain.PageService
	nodeService  domain.NodeService
	logger       *zap.SugaredLogger
}

type Option func(*Service)

func WithIndexJobRepo(indexJobRepo domain.IndexJobRepository) Option {
	return func(s *Service) {
		s.indexJobRepo = indexJobRepo
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

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) CreateJob(pageID domain.PageID, shardID domain.ShardID) error {

	if _, err := s.indexJobRepo.Get(pageID); err == nil {
		return nil, domain.ErrIndexJobAlreadyExists
	}

	job := &domain.IndexJob{
		PageID:    pageID,
		Status:    domain.IndexStatusPending,
		CreatedAt: time.Now(),
	}

	return s.indexJobRepo.Save(job)
}

func (s *Service) RunIndexProcess() error {
	jobs, err := s.indexJobRepo.ListByStatus(10, domain.IndexStatusPending)

	if err != nil {
		return err
	}

	for _, job := range jobs {
		page, err := s.pageService.Get(job.PageID)

		if err != nil {
			s.logger.Errorf("Error getting page: %v", err)
			continue
		}

		nodes, err := s.nodeService.RandomizedListGroupByShard(job.)

		if err != nil {
			s.logger.Errorf("Error listing nodes: %v", err)
			continue
		}

		for _, node := range nodes {
			s.logger.Infof("Indexing page %s on node %s", page.ID, node.ID)
		}

		job.Status = domain.IndexStatusCompleted

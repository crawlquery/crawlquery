package service

import (
	"crawlquery/api/domain"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	crawlJobRepo domain.CrawlJobRepository
	crawlLogRepo domain.CrawlLogRepository
	logger       *zap.SugaredLogger
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

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) CreateJob(pageID domain.PageID) error {

	if _, err := s.crawlJobRepo.Get(pageID); err == nil {
		return nil
	}

	cj := &domain.CrawlJob{
		PageID: pageID,
		Status: domain.CrawlStatusPending,
	}

	err := s.crawlJobRepo.Save(cj)
	if err != nil {
		return err
	}

	cl := &domain.CrawlLog{
		PageID:    pageID,
		Status:    domain.CrawlStatusPending,
		CreatedAt: time.Now(),
	}

	err = s.crawlLogRepo.Save(cl)

	if err != nil {
		return err
	}

	return nil
}

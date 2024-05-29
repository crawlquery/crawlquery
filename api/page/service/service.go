package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"

	"time"

	"go.uber.org/zap"
)

type Service struct {
	pageRepo     domain.PageRepository
	shardService domain.ShardService
	crawlService domain.CrawlService
	logger       *zap.SugaredLogger
}

type Option func(*Service)

func WithShardService(shardService domain.ShardService) func(*Service) {
	return func(s *Service) {
		s.shardService = shardService
	}
}

func WithPageRepo(pageRepo domain.PageRepository) func(*Service) {
	return func(s *Service) {
		s.pageRepo = pageRepo
	}
}

func WithCrawlService(crawlService domain.CrawlService) func(*Service) {
	return func(s *Service) {
		s.crawlService = crawlService
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

func (s *Service) Get(pageID domain.PageID) (*domain.Page, error) {
	page, err := s.pageRepo.Get(pageID)
	if err != nil {
		s.logger.Errorw("Error getting page", "error", err, "pageID", pageID)
		return nil, err
	}
	return page, nil
}

func (s *Service) Create(url domain.URL) (*domain.Page, error) {

	pageID := util.PageID(url)

	if _, err := s.pageRepo.Get(pageID); err == nil {
		s.logger.Errorw("Page already exists", "pageID", pageID)
		return nil, domain.ErrPageAlreadyExists
	}

	page := &domain.Page{
		ID:        pageID,
		URL:       url,
		CreatedAt: time.Now(),
	}

	shardID, err := s.shardService.GetURLShardID(url)

	if err != nil {
		s.logger.Errorw("Error getting shard ID", "error", err, "pageID", pageID)
		return nil, err
	}

	page.ShardID = shardID

	if err := s.pageRepo.Create(page); err != nil {
		s.logger.Errorw("Error creating page", "error", err, "pageID", pageID)
		return nil, err
	}

	if err := s.crawlService.CreateJob(page); err != nil {
		s.logger.Errorw("Error creating crawl job", "error", err, "page", page)
		return nil, err
	}

	return page, nil
}

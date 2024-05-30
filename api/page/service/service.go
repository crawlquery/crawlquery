package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"

	"time"

	"go.uber.org/zap"
)

type Service struct {
	pageRepo     domain.PageRepository
	eventService domain.EventService
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

func WithEventService(eventService domain.EventService) func(*Service) {
	return func(s *Service) {
		s.eventService = eventService
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

	s.eventService.Subscribe(domain.LinkCreatedKey, s.handleLinkCreated)
}

func (s *Service) handleLinkCreated(event domain.Event) {
	linkCreated := event.(*domain.LinkCreated)

	_, err := s.pageRepo.Get(util.PageID(linkCreated.DstURL))
	if err == domain.ErrPageNotFound {
		_, err = s.Create(linkCreated.DstURL)
		if err != nil {
			s.logger.Errorw("Error creating page", "error", err, "url", linkCreated.DstURL)
			return
		}
	}
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

	s.eventService.Publish(&domain.PageCreated{
		Page: page,
	})

	return page, nil
}

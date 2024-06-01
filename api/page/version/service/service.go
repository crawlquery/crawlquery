package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	versionRepo  domain.PageVersionRepository
	eventService domain.EventService
	logger       *zap.SugaredLogger
}

type Option func(*Service)

func WithVersionRepo(versionRepo domain.PageVersionRepository) Option {
	return func(s *Service) {
		s.versionRepo = versionRepo
	}
}

func WithEventService(eventService domain.EventService) Option {
	return func(s *Service) {
		s.eventService = eventService
	}
}

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithEventListeners() Option {
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

	s.eventService.Subscribe(domain.CrawlCompletedKey, s.onCrawlCompleted)
}

func (s *Service) onCrawlCompleted(event domain.Event) {
	crawlCompleted := event.(*domain.CrawlCompleted)

	_, err := s.Create(crawlCompleted.PageID, crawlCompleted.ContentHash)

	if err != nil {
		s.logger.Errorw("Failed to create page version", "error", err)
		return

	}
}

func (s *Service) Create(pageID domain.PageID, contentHash domain.ContentHash) (*domain.PageVersion, error) {
	pageVersion := &domain.PageVersion{
		ID:          domain.PageVersionID(util.UUIDString()),
		PageID:      pageID,
		ContentHash: contentHash,
		CreatedAt:   time.Now(),
	}

	err := s.versionRepo.Create(pageVersion)
	if err != nil {
		return nil, err
	}

	return pageVersion, nil
}

func (s *Service) Get(id domain.PageVersionID) (*domain.PageVersion, error) {
	pageVersion, err := s.versionRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return pageVersion, nil
}

func (s *Service) ListByPageID(pageID domain.PageID) ([]*domain.PageVersion, error) {
	versions, err := s.versionRepo.ListByPageID(pageID)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

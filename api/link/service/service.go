package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	eventService domain.EventService
	linkRepo     domain.LinkRepository
	logger       *zap.SugaredLogger
}

type Option func(*Service)

func WithEventService(eventService domain.EventService) Option {
	return func(s *Service) {
		s.eventService = eventService
	}
}

func WithLinkRepo(linkRepo domain.LinkRepository) Option {
	return func(s *Service) {
		s.linkRepo = linkRepo
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
	s.eventService.Subscribe(domain.CrawlCompletedKey, s.handleCrawlCompleted)
}

var trackingParams = []string{
	"utm_source", "utm_medium", "utm_platform", "utm_campaign", "utm_term", "utm_content", "gclid", "fbclid",
}

func normalizeURL(rawURL domain.URL) (domain.URL, error) {
	parsedURL, err := url.Parse(string(rawURL))
	if err != nil {
		return "", err
	}

	queryParams := parsedURL.Query()
	for _, param := range trackingParams {
		queryParams.Del(param)
	}
	parsedURL.RawQuery = queryParams.Encode()

	parsedURL.Host = strings.ToLower(parsedURL.Host)
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)

	return domain.URL(parsedURL.String()), nil
}

func (s *Service) handleCrawlCompleted(e domain.Event) {
	crawlCompletedEvent := e.(*domain.CrawlCompleted)
	for _, link := range crawlCompletedEvent.Links {
		_, err := s.Create(crawlCompletedEvent.PageID, link)
		if err != nil {
			s.logger.Errorw("Error creating link", "error", err)
		}
	}
}

func (s *Service) Create(src domain.PageID, dst domain.URL) (*domain.Link, error) {
	normalizedDst, err := normalizeURL(dst)
	if err != nil {
		return nil, err
	}
	link := &domain.Link{
		SrcID:     src,
		DstID:     util.PageID(normalizedDst),
		CreatedAt: time.Now(),
	}

	err = s.linkRepo.Create(link)

	if err != nil {
		return nil, err
	}

	s.eventService.Publish(&domain.LinkCreated{Link: link, DstURL: normalizedDst})

	return link, nil
}

func (s *Service) GetAll() ([]*domain.Link, error) {
	return s.linkRepo.GetAll()
}

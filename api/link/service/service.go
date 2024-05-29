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
	linkRepo    domain.LinkRepository
	pageService domain.PageService
	logger      *zap.SugaredLogger
}

type Option func(*Service)

func WithLinkRepo(linkRepo domain.LinkRepository) Option {
	return func(s *Service) {
		s.linkRepo = linkRepo
	}
}

func WithPageService(pageService domain.PageService) Option {
	return func(s *Service) {
		s.pageService = pageService
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

	_, err = s.pageService.Create(normalizedDst)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *Service) GetAll() ([]*domain.Link, error) {
	return s.linkRepo.GetAll()
}

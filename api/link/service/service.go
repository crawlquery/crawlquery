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
	linkRepo        domain.LinkRepository
	crawlJobService domain.CrawlJobService
	logger          *zap.SugaredLogger
}

func NewService(
	linkRepo domain.LinkRepository,
	crawlJobService domain.CrawlJobService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		linkRepo:        linkRepo,
		logger:          logger,
		crawlJobService: crawlJobService,
	}
}

var trackingParams = []string{
	"utm_source", "utm_medium", "utm_platform", "utm_campaign", "utm_term", "utm_content", "gclid", "fbclid",
}

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
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

	return parsedURL.String(), nil
}

func (s *Service) Create(src, dst string) (*domain.Link, error) {
	normalizedSrc, err := normalizeURL(src)
	if err != nil {
		return nil, err
	}

	normalizedDst, err := normalizeURL(dst)
	if err != nil {
		return nil, err
	}
	link := &domain.Link{
		SrcID:     util.PageID(normalizedSrc),
		DstID:     util.PageID(normalizedDst),
		CreatedAt: time.Now(),
	}

	err = s.linkRepo.Create(link)

	if err != nil {
		return nil, err
	}

	_, err = s.crawlJobService.Create(normalizedDst)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *Service) GetAll() ([]*domain.Link, error) {
	return s.linkRepo.GetAll()
}

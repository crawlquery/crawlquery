package service

import (
	apiDomain "crawlquery/api/domain"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	resRepo apiDomain.CrawlRestrictionRepository
	logger  *zap.SugaredLogger
}

func NewService(
	resRepo apiDomain.CrawlRestrictionRepository,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		resRepo: resRepo,
		logger:  logger,
	}
}

func (s *Service) GetRestriction(domain string) (bool, *time.Time) {
	found, err := s.resRepo.Get(domain)

	if err == apiDomain.ErrCrawlRestrictionNotFound {
		return false, nil
	}

	if found == nil {
		return false, nil
	}

	if found.Until.Time.After(time.Now()) {
		return true, &found.Until.Time
	}

	return false, nil
}

func (s *Service) Restrict(domain string) error {
	if restricted, _ := s.GetRestriction(domain); restricted {
		return apiDomain.ErrCrawlRestrictionAlreadyExists
	}

	restriction := &apiDomain.CrawlRestriction{
		Domain: domain,
		Until: sql.NullTime{
			Valid: true,
			Time:  time.Now().Add(time.Minute * 5),
		},
	}

	s.logger.Infow("CrawlRestrictionService.Restrict: setting restriction", "restriction", restriction)

	return s.resRepo.Set(restriction)
}

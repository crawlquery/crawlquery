package service

import (
	apiDomain "crawlquery/api/domain"
	"database/sql"
	"time"
)

type Service struct {
	resRepo apiDomain.CrawlRestrictionRepository
}

func NewService(resRepo apiDomain.CrawlRestrictionRepository) *Service {
	return &Service{
		resRepo: resRepo,
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

	return s.resRepo.Set(restriction)
}

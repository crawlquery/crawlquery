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

func (s *Service) HasRestriction(domain string) bool {
	found, err := s.resRepo.Get(domain)

	if err == apiDomain.ErrCrawlRestrictionNotFound {
		return false
	}

	if found == nil {
		return false
	}

	if found.Until.Time.After(time.Now()) {
		return true
	}

	return false
}

func (s *Service) Restrict(domain string) error {
	if s.HasRestriction(domain) {
		return apiDomain.ErrCrawlRestrictionAlreadyExists
	}

	restriction := &apiDomain.CrawlRestriction{
		Domain: domain,
		Until: sql.NullTime{
			Valid: true,
			Time:  time.Now().Add(time.Hour),
		},
	}

	return s.resRepo.Set(restriction)
}

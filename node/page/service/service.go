package service

import (
	"crawlquery/node/domain"
	sharedDomain "crawlquery/pkg/domain"
)

type Service struct {
	pageRepo domain.PageRepository
}

func NewService(pr domain.PageRepository) *Service {
	return &Service{
		pageRepo: pr,
	}
}

func (s *Service) Save(pageID string, url string) (*sharedDomain.Page, error) {

	page := &sharedDomain.Page{
		ID:  pageID,
		URL: url,
	}

	s.pageRepo.Save(pageID, page)

	return page, nil
}

func (s *Service) Get(pageID string) (*sharedDomain.Page, error) {
	return s.pageRepo.Get(pageID)
}

package service

import "crawlquery/node/domain"

type Service struct {
	repo domain.HTMLRepository
}

func NewService(repo domain.HTMLRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Save(pageID string, html []byte) error {
	return s.repo.Save(pageID, html)
}

func (s *Service) Get(pageID string) ([]byte, error) {
	return s.repo.Get(pageID)
}

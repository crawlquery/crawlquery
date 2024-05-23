package service

import "crawlquery/node/domain"

type Service struct {
	pageService domain.PageService
}

func NewService(
	pageService domain.PageService,
) *Service {
	return &Service{
		pageService: pageService,
	}
}

func (s *Service) Page() ([]byte, error) {
	return s.pageService.JSON()
}

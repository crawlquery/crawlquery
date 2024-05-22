package service

import "crawlquery/node/domain"

type Service struct {
	pageService    domain.PageService
	keywordService domain.KeywordService
}

func NewService(
	pageService domain.PageService,
	keywordService domain.KeywordService,
) *Service {
	return &Service{
		pageService:    pageService,
		keywordService: keywordService,
	}
}

func (s *Service) Page() ([]byte, error) {
	return s.pageService.JSON()
}

func (s *Service) Keyword() ([]byte, error) {
	return s.keywordService.JSON()
}

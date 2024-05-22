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

func (s *Service) PageDump() ([]byte, error) {
	return s.pageService.JSON()
}

func (s *Service) KeywordDump() ([]byte, error) {
	return s.keywordService.JSON()
}

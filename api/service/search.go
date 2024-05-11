package service

import (
	"crawlquery/pkg/domain"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) Search(term string) []domain.Result {
	return []domain.Result{}
}

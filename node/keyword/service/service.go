package service

import (
	"crawlquery/node/domain"
	"strings"
)

type Service struct {
	repo domain.KeywordRepository
}

func NewService(repo domain.KeywordRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) UpdatePageKeywords(pageID string, keywordsSplit [][]string) error {

	var keywords []string

	for _, keyword := range keywordsSplit {
		keywords = append(keywords, strings.Join(keyword, " "))
	}

	err := s.repo.RemovePageKeywords(pageID)

	if err != nil {
		return err
	}

	return s.repo.AddPageKeywords(pageID, keywords)
}

func (s *Service) GetPageIDsByKeyword(keyword string) ([]string, error) {
	return s.repo.GetPages(keyword)
}

package service

import (
	"crawlquery/node/domain"
)

type Service struct {
	repo domain.KeywordOccurrenceRepository
}

func NewService(repo domain.KeywordOccurrenceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Count() (int, error) {
	return s.repo.Count()
}

func (s *Service) GetKeywordMatches(keywords []domain.Keyword) ([]domain.KeywordMatch, error) {
	var matches []domain.KeywordMatch

	for _, keyword := range keywords {
		occurrences, err := s.repo.GetAll(keyword)
		if err != nil {
			if err == domain.ErrKeywordNotFound {
				continue
			}
			return nil, err
		}

		matches = append(matches, domain.KeywordMatch{
			Keyword:     keyword,
			Occurrences: occurrences,
		})
	}

	return matches, nil
}

func (s *Service) UpdateOccurrences(pageID string, keywordOccurrences map[domain.Keyword]domain.KeywordOccurrence) error {
	for keyword, occurrence := range keywordOccurrences {
		err := s.repo.Add(keyword, occurrence)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) RemoveOccurrencesForPageID(pageID string) error {
	return s.repo.RemoveForPageID(pageID)
}

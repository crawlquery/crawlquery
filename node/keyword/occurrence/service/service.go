package service

import (
	"crawlquery/node/domain"
)

type KeywordOccurrenceService struct {
	repo domain.KeywordOccurrenceRepository
}

func NewKeywordOccurrenceService(repo domain.KeywordOccurrenceRepository) *KeywordOccurrenceService {
	return &KeywordOccurrenceService{repo: repo}
}

func (s *KeywordOccurrenceService) GetKeywordMatches(keywords []domain.Keyword) ([]domain.KeywordMatch, error) {
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

func (s *KeywordOccurrenceService) Update(pageID string, keywordOccurrences map[domain.Keyword]domain.KeywordOccurrence) error {
	for keyword, occurrence := range keywordOccurrences {
		err := s.repo.Add(keyword, occurrence)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *KeywordOccurrenceService) RemoveForPageID(pageID string) error {
	return s.repo.RemoveForPageID(pageID)
}

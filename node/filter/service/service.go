package service

import "crawlquery/node/domain"

type Service struct {
	explicitWords []string
}

func NewService() *Service {
	return &Service{
		explicitWords: []string{
			"porn",
			"fuck",
			"sex",
		},
	}
}

func (s *Service) IsExplicitWord(word string) bool {
	for _, w := range s.explicitWords {
		if w == word {
			return true
		}
	}

	return false
}

func (s *Service) PostingsContainsExplicitWords(postings map[string]*domain.Posting) bool {
	for keyword := range postings {
		if s.IsExplicitWord(keyword) {
			return true
		}
	}

	return false
}

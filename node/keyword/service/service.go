package service

import "crawlquery/node/domain"

type Service struct {
	repo domain.KeywordRepository
}

func NewService(repo domain.KeywordRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SavePostings(keywordPostings map[string]*domain.Posting) error {
	for keyword, posting := range keywordPostings {
		if err := s.repo.SavePosting(keyword, posting); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetPostings(keyword string) ([]*domain.Posting, error) {
	postings, err := s.repo.GetPostings(keyword)
	if err != nil {
		return nil, err
	}

	return postings, nil
}

func (s *Service) FuzzySearch(token string) ([]string, error) {
	results := s.repo.FuzzySearch(token)
	return results, nil
}

func (s *Service) RemovePostingsByPageID(pageID string) error {
	return s.repo.RemovePostingsByPageID(pageID)
}

package mem

import (
	"crawlquery/node/domain"
	"strings"
)

type Repository struct {
	keywordPostings map[string][]*domain.Posting
}

func NewRepository() *Repository {
	return &Repository{
		keywordPostings: make(map[string][]*domain.Posting),
	}
}

func (r *Repository) SavePosting(token string, posting *domain.Posting) error {

	_, ok := r.keywordPostings[token]

	if !ok {
		r.keywordPostings[token] = make([]*domain.Posting, 0)
	}

	r.keywordPostings[token] = append(r.keywordPostings[token], posting)
	return nil
}

func (r *Repository) FuzzySearch(token string) []string {
	results := []string{}

	for tokens, _ := range r.keywordPostings {
		if strings.Contains(tokens, token) {
			results = append(results, tokens)
		}
	}
	return results
}

func (r *Repository) GetPostings(token string) ([]*domain.Posting, error) {
	postings, ok := r.keywordPostings[token]
	if !ok {
		return nil, nil
	}
	return postings, nil
}

func (r *Repository) RemovePostingsByPageID(pageID string) error {
	for token, postings := range r.keywordPostings {
		for i, posting := range postings {
			if posting.PageID == pageID {
				r.keywordPostings[token] = append(r.keywordPostings[token][:i], r.keywordPostings[token][i+1:]...)
			}
		}
	}
	return nil
}

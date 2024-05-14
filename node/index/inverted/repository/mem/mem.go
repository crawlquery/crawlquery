package mem

import (
	"crawlquery/node/domain"
	"strings"
)

type Repository struct {
	forwardIndex map[string][]*domain.Posting
}

func NewRepository() *Repository {
	return &Repository{
		forwardIndex: make(map[string][]*domain.Posting),
	}
}

func (r *Repository) Save(token string, posting *domain.Posting) error {

	_, ok := r.forwardIndex[token]

	if !ok {
		r.forwardIndex[token] = make([]*domain.Posting, 0)
	}

	r.forwardIndex[token] = append(r.forwardIndex[token], posting)
	return nil
}

func (r *Repository) FuzzySearch(token string) map[string]float64 {
	results := make(map[string]float64)

	for key, postings := range r.forwardIndex {
		// Check if the term is a substring of the key
		if strings.Contains(key, token) {
			for _, posting := range postings {
				// Add or increase the score based on frequency
				results[posting.PageID] += float64(posting.Frequency)
			}
		}
	}
	return results
}

func (r *Repository) Get(pageID string) ([]*domain.Posting, error) {
	postings, ok := r.forwardIndex[pageID]
	if !ok {
		return nil, nil
	}
	return postings, nil
}

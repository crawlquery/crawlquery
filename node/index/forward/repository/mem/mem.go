package mem

import "crawlquery/pkg/domain"

type Repository struct {
	forwardIndex map[string]*domain.Page
}

func NewRepository() *Repository {
	return &Repository{
		forwardIndex: make(map[string]*domain.Page),
	}
}

func (r *Repository) Save(pageID string, page *domain.Page) error {
	r.forwardIndex[pageID] = page
	return nil
}

func (r *Repository) Get(pageID string) (*domain.Page, error) {
	page, ok := r.forwardIndex[pageID]
	if !ok {
		return nil, nil
	}
	return page, nil
}

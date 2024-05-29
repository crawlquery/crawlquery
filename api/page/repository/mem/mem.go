package mem

import "crawlquery/api/domain"

type Repository struct {
	pages map[domain.PageID]*domain.Page
}

func NewRepository() *Repository {
	return &Repository{
		pages: make(map[domain.PageID]*domain.Page),
	}
}

func (r *Repository) Get(id domain.PageID) (*domain.Page, error) {
	page, ok := r.pages[id]
	if !ok {
		return nil, domain.ErrPageNotFound
	}
	return page, nil
}

func (r *Repository) Create(p *domain.Page) error {
	r.pages[p.ID] = p
	return nil
}

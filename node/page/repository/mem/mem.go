package mem

import (
	"crawlquery/node/domain"
	sharedDomain "crawlquery/pkg/domain"
)

type Repository struct {
	pages map[string]*sharedDomain.Page
}

func NewRepository() *Repository {
	return &Repository{
		pages: make(map[string]*sharedDomain.Page),
	}
}

func (r *Repository) Save(pageID string, page *sharedDomain.Page) error {
	r.pages[pageID] = page
	return nil
}

func (r *Repository) Get(pageID string) (*sharedDomain.Page, error) {
	page, ok := r.pages[pageID]
	if !ok {
		return nil, domain.ErrPageNotFound
	}
	return page, nil
}

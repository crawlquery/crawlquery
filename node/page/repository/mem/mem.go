package mem

import (
	"crawlquery/node/domain"
)

type Repository struct {
	pages      map[string]*domain.Page
	pageHashes map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		pages:      make(map[string]*domain.Page),
		pageHashes: make(map[string]string),
	}
}

func (r *Repository) Save(pageID string, page *domain.Page) error {
	r.pages[pageID] = page
	return nil
}

func (r *Repository) Delete(pageID string) error {
	delete(r.pages, pageID)
	return nil
}

func (r *Repository) Get(pageID string) (*domain.Page, error) {
	page, ok := r.pages[pageID]
	if !ok {
		return nil, domain.ErrPageNotFound
	}
	return page, nil
}

func (r *Repository) Count() (int, error) {
	return len(r.pages), nil
}

func (r *Repository) GetByIDs(pageIDs []string) (map[string]*domain.Page, error) {
	pages := make(map[string]*domain.Page)
	for _, pageID := range pageIDs {
		page, ok := r.pages[pageID]
		if ok {
			pages[pageID] = page
		}
	}
	return pages, nil
}

func (r *Repository) GetAll() (map[string]*domain.Page, error) {
	return r.pages, nil
}

func (r *Repository) UpdateHash(pageID string, hash string) error {
	r.pageHashes[pageID] = hash
	return nil
}

func (r *Repository) DeleteHash(pageID string) error {
	delete(r.pageHashes, pageID)
	return nil
}

func (r *Repository) GetHash(pageID string) (string, error) {
	hash, ok := r.pageHashes[pageID]
	if !ok {
		return "", domain.ErrHashNotFound
	}
	return hash, nil
}

func (r *Repository) GetHashes() (map[string]string, error) {
	return r.pageHashes, nil
}

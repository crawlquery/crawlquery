package mem

import (
	"crawlquery/node/domain"
	sharedDomain "crawlquery/pkg/domain"
)

type Repository struct {
	pages      map[string]*sharedDomain.Page
	pageHashes map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		pages:      make(map[string]*sharedDomain.Page),
		pageHashes: make(map[string]string),
	}
}

func (r *Repository) Save(pageID string, page *sharedDomain.Page) error {
	r.pages[pageID] = page
	return nil
}

func (r *Repository) Delete(pageID string) error {
	delete(r.pages, pageID)
	return nil
}

func (r *Repository) Get(pageID string) (*sharedDomain.Page, error) {
	page, ok := r.pages[pageID]
	if !ok {
		return nil, domain.ErrPageNotFound
	}
	return page, nil
}

func (r *Repository) GetAll() (map[string]*sharedDomain.Page, error) {
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

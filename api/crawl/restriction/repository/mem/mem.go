package mem

import (
	apidomain "crawlquery/api/domain"
	"sync"
)

type Repository struct {
	restrictions map[string]*apidomain.CrawlRestriction

	mutex sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		restrictions: make(map[string]*apidomain.CrawlRestriction),
		mutex:        sync.Mutex{},
	}
}

func (r *Repository) Get(domain string) (*apidomain.CrawlRestriction, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	restriction, ok := r.restrictions[domain]

	if !ok {
		return nil, apidomain.ErrCrawlRestrictionNotFound
	}

	return restriction, nil
}

func (r *Repository) Set(res *apidomain.CrawlRestriction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.restrictions[res.Domain] = res
	return nil
}

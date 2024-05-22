package mem

import (
	apidomain "crawlquery/api/domain"
	"crawlquery/pkg/util"
	"database/sql"
	"sync"
	"time"
)

type Repository struct {
	locks map[string]*apidomain.Lock

	mutex sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		locks: make(map[string]*apidomain.Lock),
		mutex: sync.Mutex{},
	}
}

func (r *Repository) IsLocked(domain string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, ok := r.locks[domain]
	return ok
}

func (r *Repository) Lock(domain string) (string, error) {

	if r.IsLocked(domain) {
		return "", apidomain.ErrDomainLocked
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.locks[domain] = &apidomain.Lock{
		Domain: domain,
		Key:    util.UUID(),
		LockedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	return r.locks[domain].Key, nil
}

func (r *Repository) Unlock(domain string, key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if lock, ok := r.locks[domain]; !ok || lock.Key != key {
		return apidomain.ErrInvalidLockKey
	}

	delete(r.locks, domain)
	return nil
}

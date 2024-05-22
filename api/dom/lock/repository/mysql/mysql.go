package mysql

import (
	"crawlquery/pkg/util"
	"database/sql"
	"time"

	apidomain "crawlquery/api/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) IsLocked(domain string) bool {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM domain_locks WHERE domain = ? AND locked_at IS NOT NULL", domain).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *Repository) Lock(domain string) (string, error) {
	if r.IsLocked(domain) {
		return "", apidomain.ErrDomainLocked
	}

	key := util.UUID()

	_, err := r.db.Exec("INSERT INTO domain_locks (domain, `key`, locked_at) VALUES (?, ?, ?)", domain, key, time.Now())

	if err != nil {
		return "", err
	}

	return key, nil
}

func (r *Repository) Unlock(domain string, key string) error {

	if !r.IsLocked(domain) {
		return apidomain.ErrDomainNotLocked
	}

	var keyCheck string
	err := r.db.QueryRow("SELECT `key` FROM domain_locks WHERE domain = ? AND `key` = ?", domain, key).Scan(&keyCheck)

	if err != nil {
		return apidomain.ErrInvalidLockKey
	}

	_, err = r.db.Exec("DELETE FROM domain_locks WHERE domain = ?", domain)

	if err != nil {
		return err
	}
	return nil
}

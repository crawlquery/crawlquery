package mysql

import (
	"database/sql"

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

func (r *Repository) Get(domain string) (*apidomain.CrawlRestriction, error) {
	var restriction apidomain.CrawlRestriction
	err := r.db.QueryRow("SELECT domain, until FROM crawl_restrictions WHERE domain = ? AND until IS NOT NULL", domain).Scan(&restriction.Domain, &restriction.Until)

	if err != nil {
		return nil, apidomain.ErrCrawlRestrictionNotFound
	}

	return &restriction, nil
}

func (r *Repository) Set(res *apidomain.CrawlRestriction) error {
	_, err := r.db.Exec("INSERT INTO crawl_restrictions (domain, until) VALUES (?, ?) ON DUPLICATE KEY UPDATE until = ?", res.Domain, res.Until, res.Until)

	if err != nil {
		return err
	}

	return nil
}

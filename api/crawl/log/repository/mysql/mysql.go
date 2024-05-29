package mysql

import (
	"crawlquery/api/domain"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(cl *domain.CrawlLog) error {
	_, err := r.db.Exec("INSERT INTO crawl_logs (id, page_id, status, info, created_at) VALUES (?, ?, ?, ?, ?)", cl.ID, cl.PageID, cl.Status, cl.Info, cl.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

package mysql

import (
	"crawlquery/api/domain"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Get(pageID domain.PageID) (float64, error) {
	var rank float64
	err := r.db.QueryRow("SELECT `rank` FROM page_ranks WHERE page_id = ?", pageID).Scan(&rank)
	if err != nil {
		return 0, err
	}
	return rank, nil
}

func (r *Repository) Update(pageID domain.PageID, rank float64, createdAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO page_ranks (page_id, `rank`, created_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `rank` = ?", pageID, rank, createdAt, pageID)
	if err != nil {
		return err
	}
	return nil
}

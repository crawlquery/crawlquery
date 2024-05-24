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

func (r *Repository) Get(id string) (*domain.Page, error) {
	var page domain.Page

	err := r.db.QueryRow("SELECT id, shard_id, hash, created_at FROM pages WHERE id = ?", id).Scan(&page.ID, &page.ShardID, &page.Hash, &page.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPageNotFound
		}
		return nil, err
	}

	return &page, nil
}

func (r *Repository) Create(p *domain.Page) error {
	_, err := r.db.Exec("INSERT INTO pages (id, shard_id, hash, created_at) VALUES (?, ?, ?, ?)", p.ID, p.ShardID, p.Hash, p.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
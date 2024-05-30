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

func (r *Repository) Get(id domain.PageVersionID) (*domain.PageVersion, error) {
	var version domain.PageVersion
	err := r.db.QueryRow("SELECT id, page_id, content_hash, created_at FROM page_versions WHERE id = ?", id).Scan(&version.ID, &version.PageID, &version.ContentHash, &version.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPageVersionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *Repository) Create(v *domain.PageVersion) error {
	_, err := r.db.Exec("INSERT INTO page_versions (id, page_id, content_hash, created_at) VALUES (?, ?, ?, ?)", v.ID, v.PageID, v.ContentHash, v.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ListByPageID(pageID domain.PageID) ([]*domain.PageVersion, error) {
	rows, err := r.db.Query("SELECT id, page_id, content_hash, created_at FROM page_versions WHERE page_id = ?", pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*domain.PageVersion
	for rows.Next() {
		var version domain.PageVersion
		err := rows.Scan(&version.ID, &version.PageID, &version.ContentHash, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, &version)
	}
	return versions, nil
}

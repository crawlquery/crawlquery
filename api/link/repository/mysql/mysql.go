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

func (r *Repository) Create(link *domain.Link) error {
	_, err := r.db.Exec("INSERT INTO links (src_id, dst_id, created_at) VALUES (?, ?, ?)", link.SrcID, link.DstID, link.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAll() ([]*domain.Link, error) {
	rows, err := r.db.Query("SELECT src_id, dst_id, created_at FROM links")

	if err != nil {
		return nil, err
	}

	var links []*domain.Link
	for rows.Next() {
		var link domain.Link
		err = rows.Scan(&link.SrcID, &link.DstID, &link.CreatedAt)
		if err != nil {
			return nil, err
		}

		links = append(links, &link)
	}

	return links, nil
}

func (r *Repository) GetAllBySrcID(srcID domain.PageID) ([]*domain.Link, error) {
	rows, err := r.db.Query("SELECT src_id, dst_id, created_at FROM links WHERE src_id = ?", srcID)
	if err != nil {
		return nil, err
	}

	var links []*domain.Link
	for rows.Next() {
		var link domain.Link
		err = rows.Scan(&link.SrcID, &link.DstID, &link.CreatedAt)
		if err != nil {
			return nil, err
		}

		links = append(links, &link)
	}

	return links, nil
}

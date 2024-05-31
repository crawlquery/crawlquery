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

func (r *Repository) Save(log *domain.IndexLog) error {
	_, err := r.db.Exec("INSERT INTO index_logs (id, page_id, status, info, created_at) VALUES (?, ?, ?, ?, ?)", log.ID, log.PageID, log.Status, log.Info, log.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ListByPageID(pageID domain.PageID) ([]*domain.IndexLog, error) {
	rows, err := r.db.Query("SELECT id, page_id, status, info, created_at FROM index_logs WHERE page_id = ?", pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.IndexLog
	for rows.Next() {
		var log domain.IndexLog
		err := rows.Scan(&log.ID, &log.PageID, &log.Status, &log.Info, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, nil
}

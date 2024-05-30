package mysql

import (
	"crawlquery/api/domain"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(job *domain.IndexJob) error {
	_, err := r.db.Exec("INSERT INTO index_jobs (page_id, shard_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = ?, updated_at = ?", job.PageID, job.ShardID, job.Status, job.CreatedAt, job.UpdatedAt, job.Status, job.UpdatedAt)

	return err
}

func (r *Repository) Get(pageID domain.PageID) (*domain.IndexJob, error) {
	var job domain.IndexJob
	err := r.db.QueryRow("SELECT page_id, shard_id, status, created_at, updated_at FROM index_jobs WHERE page_id = ?", pageID).Scan(&job.PageID, &job.ShardID, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrIndexJobNotFound
		}
		return nil, err
	}
	return &job, nil
}

func (r *Repository) ListByStatus(limit int, status domain.IndexStatus) ([]*domain.IndexJob, error) {
	rows, err := r.db.Query("SELECT page_id, shard_id, status, created_at, updated_at FROM index_jobs WHERE status = ? LIMIT ?", status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.IndexJob
	for rows.Next() {
		var job domain.IndexJob
		err := rows.Scan(&job.PageID, &job.ShardID, &job.Status, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

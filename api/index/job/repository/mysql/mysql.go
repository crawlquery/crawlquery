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

func (r *Repository) Create(job *domain.IndexJob) (*domain.IndexJob, error) {
	_, err := r.db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		job.ID, job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.CreatedAt)

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (r *Repository) Get(id string) (*domain.IndexJob, error) {
	var job domain.IndexJob

	err := r.db.QueryRow("SELECT id, page_id, backoff_until, last_indexed_at, failed_reason, created_at FROM index_jobs WHERE id = ?", id).
		Scan(&job.ID, &job.PageID, &job.BackoffUntil, &job.LastIndexedAt, &job.FailedReason, &job.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrIndexJobNotFound
		}
		return nil, err
	}

	return &job, nil
}

func (r *Repository) Update(job *domain.IndexJob) error {
	_, err := r.db.Exec("UPDATE index_jobs SET page_id = ?, backoff_until = ?, last_indexed_at = ?, failed_reason = ? WHERE id = ?",
		job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Next() (*domain.IndexJob, error) {
	var job domain.IndexJob

	err := r.db.QueryRow("SELECT id, page_id, backoff_until, last_indexed_at, failed_reason, created_at FROM index_jobs WHERE last_indexed_at IS NULL and (backoff_until IS NULL OR backoff_until < NOW()) ORDER BY created_at ASC LIMIT 1").
		Scan(&job.ID, &job.PageID, &job.BackoffUntil, &job.LastIndexedAt, &job.FailedReason, &job.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrIndexJobNotFound
		}
		return nil, err
	}

	return &job, nil
}

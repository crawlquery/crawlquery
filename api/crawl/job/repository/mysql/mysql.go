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
	return &Repository{db}
}

func (r *Repository) Create(j *domain.CrawlJob) error {
	_, err := r.db.Exec("INSERT INTO crawl_jobs (id, url, url_hash, created_at) VALUES (?, ?, ?, ?)", j.ID, j.URL, j.URLHash, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(j *domain.CrawlJob) error {
	res, err := r.db.Exec("UPDATE crawl_jobs SET backoff_until = ?, last_crawled_at = ?, failed_reason = ? WHERE id = ?", j.BackoffUntil, j.LastCrawledAt, j.FailedReason, j.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNoRowsUpdated
	}

	return nil
}

func (r *Repository) Get(id string) (*domain.CrawlJob, error) {
	row := r.db.QueryRow("SELECT id, url, url_hash, backoff_until, failed_reason, last_crawled_at, created_at FROM crawl_jobs WHERE id = ?", id)

	var job domain.CrawlJob
	err := row.Scan(&job.ID, &job.URL, &job.URLHash, &job.BackoffUntil, &job.FailedReason, &job.LastCrawledAt, &job.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrCrawlJobNotFound
	}

	return &job, err
}

func (r *Repository) First() (*domain.CrawlJob, error) {
	row := r.db.QueryRow("SELECT id, url, url_hash, backoff_until, failed_reason, last_crawled_at, created_at FROM crawl_jobs ORDER BY created_at ASC LIMIT 1")

	var job domain.CrawlJob
	err := row.Scan(&job.ID, &job.URL, &job.URLHash, &job.BackoffUntil, &job.FailedReason, &job.LastCrawledAt, &job.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCrawlJobNotFound
	}

	return &job, err
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM crawl_jobs WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

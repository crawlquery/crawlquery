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
	_, err := r.db.Exec("INSERT INTO crawl_jobs (id, url, created_at) VALUES (?, ?, ?)", j.ID, j.URL, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(id string) (*domain.CrawlJob, error) {
	row := r.db.QueryRow("SELECT id, url, created_at FROM crawl_jobs WHERE id = ?", id)

	var job domain.CrawlJob
	err := row.Scan(&job.ID, &job.URL, &job.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrCrawlJobNotFound
	}

	return &job, err
}

func (r *Repository) First() (*domain.CrawlJob, error) {
	row := r.db.QueryRow("SELECT id, url, created_at FROM crawl_jobs ORDER BY created_at ASC LIMIT 1")

	var job domain.CrawlJob
	err := row.Scan(&job.ID, &job.URL, &job.CreatedAt)
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

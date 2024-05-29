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

func (r *Repository) Get(pageID domain.PageID) (*domain.CrawlJob, error) {
	var job domain.CrawlJob
	err := r.db.QueryRow("SELECT page_id, status, created_at, updated_at FROM crawl_jobs WHERE page_id = ?", pageID).Scan(&job.PageID, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCrawlJobNotFound
		}
		return nil, err
	}
	return &job, nil
}

func (r *Repository) Save(job *domain.CrawlJob) error {
	_, err := r.db.Exec("INSERT INTO crawl_jobs (page_id, status, created_at, updated_at) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = ?, updated_at = ?", job.PageID, job.Status, job.CreatedAt, job.UpdatedAt, job.Status, job.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

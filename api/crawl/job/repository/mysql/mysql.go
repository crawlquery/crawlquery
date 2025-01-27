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
	err := r.db.QueryRow("SELECT page_id, url, shard_id, status, created_at, updated_at FROM crawl_jobs WHERE page_id = ?", pageID).Scan(&job.PageID, &job.URL, &job.ShardID, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCrawlJobNotFound
		}
		return nil, err
	}
	return &job, nil
}

func (r *Repository) Save(job *domain.CrawlJob) error {
	_, err := r.db.Exec("INSERT INTO crawl_jobs (page_id, url, shard_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = ?, updated_at = ?", job.PageID, job.URL, job.ShardID, job.Status, job.CreatedAt, job.UpdatedAt, job.Status, job.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ListByStatus(limit int, status domain.CrawlStatus) ([]*domain.CrawlJob, error) {
	rows, err := r.db.Query("SELECT page_id, url, shard_id, status, created_at, updated_at FROM crawl_jobs WHERE status = ? LIMIT ?", status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.CrawlJob
	for rows.Next() {
		var job domain.CrawlJob
		err := rows.Scan(&job.PageID, &job.URL, &job.ShardID, &job.Status, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

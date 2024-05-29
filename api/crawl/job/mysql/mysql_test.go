package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	crawlJobRepo "crawlquery/api/crawl/job/mysql"
)

func TestGet(t *testing.T) {
	t.Run("can get a job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		crawlJobRepo := crawlJobRepo.NewRepository(db)
		crawlJob := &domain.CrawlJob{
			PageID:    "page1",
			Status:    domain.CrawlJobStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := db.Exec("INSERT INTO crawl_jobs (page_id, status, created_at, updated_at) VALUES (?, ?, ?, ?)", crawlJob.PageID, crawlJob.Status, crawlJob.CreatedAt, crawlJob.UpdatedAt)
		defer db.Exec("DELETE FROM crawl_jobs WHERE page_id = ?", crawlJob.PageID)

		if err != nil {
			t.Fatal(err)
		}

		job, err := crawlJobRepo.Get("page1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if job.PageID != "page1" {
			t.Errorf("expected page1, got %v", job.PageID)
		}

		if job.Status != domain.CrawlJobStatusPending {
			t.Errorf("expected pending, got %v", job.Status)
		}
	})

	t.Run("should return ErrCrawlJobNotFound if job does not exist", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		crawlJobRepo := crawlJobRepo.NewRepository(db)

		job, err := crawlJobRepo.Get("page1")

		if err != domain.ErrCrawlJobNotFound {
			t.Errorf("expected ErrCrawlJobNotFound, got %v", err)
		}

		if job != nil {
			t.Errorf("expected nil, got %v", job)
		}
	})
}

func TestSave(t *testing.T) {
	t.Run("can create a job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		crawlJobRepo := crawlJobRepo.NewRepository(db)
		now := time.Now()
		crawlJob := &domain.CrawlJob{
			PageID:    "page1",
			Status:    domain.CrawlJobStatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := crawlJobRepo.Save(crawlJob)
		defer db.Exec("DELETE FROM crawl_jobs WHERE page_id = ?", crawlJob.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		job, err := crawlJobRepo.Get("page1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if job.PageID != "page1" {
			t.Errorf("expected page1, got %v", job.PageID)
		}

		if job.Status != domain.CrawlJobStatusPending {
			t.Errorf("expected pending, got %v", job.Status)
		}

		if job.CreatedAt.UTC().Round(time.Second) != now.UTC().Round(time.Second) {
			t.Errorf("expected %v, got %v", now, job.CreatedAt)
		}

		if job.UpdatedAt.UTC().Round(time.Second) != now.UTC().Round(time.Second) {
			t.Errorf("expected %v, got %v", now, job.UpdatedAt)
		}
	})

	t.Run("can update a job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		crawlJobRepo := crawlJobRepo.NewRepository(db)
		now := time.Now()
		crawlJob := &domain.CrawlJob{
			PageID:    "page1",
			Status:    domain.CrawlJobStatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err := db.Exec("INSERT INTO crawl_jobs (page_id, status, created_at, updated_at) VALUES (?, ?, ?, ?)", crawlJob.PageID, crawlJob.Status, crawlJob.CreatedAt, crawlJob.UpdatedAt)
		if err != nil {
			t.Fatal(err)
		}

		err = crawlJobRepo.Save(crawlJob)

		defer db.Exec("DELETE FROM crawl_jobs WHERE page_id = ?", crawlJob.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		crawlJob.Status = domain.CrawlJobStatusInProgress
		crawlJob.UpdatedAt = time.Now()

		err = crawlJobRepo.Save(crawlJob)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		var job domain.CrawlJob

		err = db.QueryRow("SELECT page_id, status, created_at, updated_at FROM crawl_jobs WHERE page_id = ?", crawlJob.PageID).Scan(&job.PageID, &job.Status, &job.CreatedAt, &job.UpdatedAt)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if job.PageID != "page1" {
			t.Errorf("expected page1, got %v", job.PageID)
		}

		if job.Status != domain.CrawlJobStatusInProgress {
			t.Errorf("expected in progress, got %v", job.Status)
		}

		if job.UpdatedAt.UTC().Round(time.Second) != crawlJob.UpdatedAt.UTC().Round(time.Second) {
			t.Errorf("expected %v, got %v", crawlJob.UpdatedAt, job.UpdatedAt)
		}
	})
}

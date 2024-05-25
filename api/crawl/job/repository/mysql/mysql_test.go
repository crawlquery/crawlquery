package mysql_test

import (
	"crawlquery/api/crawl/job/repository/mysql"
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestCreate(t *testing.T) {
	t.Run("can create a job", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			CreatedAt: time.Now().UTC(),
		}
		err := repo.Create(job)

		// Assert
		if err != nil {
			t.Errorf("Error creating job: %v", err)
		}
		res, err := db.Query("SELECT id, url, page_id, created_at FROM crawl_jobs WHERE id = ?", job.ID)

		if err != nil {
			t.Errorf("Error querying for job: %v", err)
		}

		var id string
		var url string
		var urlHash string
		var createdAt time.Time

		for res.Next() {
			err = res.Scan(&id, &url, &urlHash, &createdAt)
			if err != nil {
				t.Errorf("Error scanning job: %v", err)
			}
		}

		if id != job.ID {
			t.Errorf("Expected ID to be %s, got %s", job.ID, id)
		}

		if url != job.URL {
			t.Errorf("Expected URL to be %s, got %s", job.URL, url)
		}

		if createdAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(createdAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, createdAt)
		}

		// Clean up
		db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)
	})

	t.Run("cant create a job with the same ID", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			PageID:    "hash",
			CreatedAt: time.Now().UTC(),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job.ID, job.PageID, job.URL, job.CreatedAt)

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)

		// Act
		err := repo.Create(job)

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("can update a job", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			PageID:    "hash",
			CreatedAt: time.Now().UTC(),
		}
		_, err := db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting job: %v", err)
		}

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)

		// Act
		job.BackoffUntil = sql.NullTime{Time: time.Now(), Valid: true}
		reason := "returned 404"
		job.FailedReason = sql.NullString{String: reason, Valid: true}
		job.LastCrawledAt = sql.NullTime{Time: time.Now(), Valid: true}

		err = repo.Update(job)

		// Assert
		if err != nil && err != domain.ErrNoRowsUpdated {
			t.Errorf("Error updating job: %v", err)
		}

		res, err := db.Query("SELECT id, url, page_id, backoff_until, failed_reason, last_crawled_at, created_at FROM crawl_jobs WHERE id = ?", job.ID)

		if err != nil {
			t.Errorf("Error querying for job: %v", err)
		}

		var id string
		var url string
		var urlHash string
		var backoffUntil sql.NullTime
		var failedReason sql.NullString
		var lastCrawledAt sql.NullTime
		var createdAt time.Time

		for res.Next() {
			err = res.Scan(&id, &url, &urlHash, &backoffUntil, &failedReason, &lastCrawledAt, &createdAt)
			if err != nil {
				t.Errorf("Error scanning job: %v", err)
			}
		}

		if id != job.ID {
			t.Errorf("Expected ID to be %s, got %s", job.ID, id)
		}

		if url != job.URL {
			t.Errorf("Expected URL to be %s, got %s", job.URL, url)
		}

		if urlHash != job.PageID {
			t.Errorf("Expected PageID to be %s, got %s", job.PageID, urlHash)
		}

		if backoffUntil.Time.Sub(job.BackoffUntil.Time) > time.Second || job.BackoffUntil.Time.Sub(backoffUntil.Time) > time.Second {
			t.Errorf("Expected BackoffUntil to be within one second of %v, got %v", job.BackoffUntil, backoffUntil)
		}

		if failedReason != job.FailedReason {
			t.Errorf("Expected FailedReason to be %s, got %s", job.FailedReason.String, failedReason.String)
		}

		if createdAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(createdAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, createdAt)
		}
	})

	t.Run("returns error if job does not exist", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			PageID:    "hash",
			CreatedAt: time.Now().UTC(),
		}

		// Act
		err := repo.Update(job)

		// Assert
		if err != domain.ErrNoRowsUpdated {
			t.Errorf("Expected error, got %v", err)
		}

	})
}

func TestGet(t *testing.T) {
	t.Run("can get a job", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			PageID:    "hash",
			CreatedAt: time.Now().UTC(),
		}
		_, err := db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting job: %v", err)
		}

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)

		// Act
		res, err := repo.Get(job.ID)

		// Assert
		if err != nil {
			t.Errorf("Error getting job: %v", err)
		}

		if res.ID != job.ID {
			t.Errorf("Expected ID to be %s, got %s", job.ID, res.ID)
		}

		if res.URL != job.URL {
			t.Errorf("Expected URL to be %s, got %s", job.URL, res.URL)
		}

		if res.CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, res.CreatedAt)
		}
	})

	t.Run("returns error if job does not exist", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		res, err := repo.Get("nonexistent")

		// Assert
		if err != domain.ErrCrawlJobNotFound {
			t.Errorf("Expected error, got %v", err)
		}

		if res != nil {
			t.Errorf("Expected nil, got %v", res)
		}
	})
}

func TestGetByPageID(t *testing.T) {
	t.Run("can get a job by page ID", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			PageID:    "hash",
			CreatedAt: time.Now().UTC(),
		}
		_, err := db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting job: %v", err)
		}

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)

		// Act
		res, err := repo.GetByPageID(job.PageID)

		// Assert
		if err != nil {
			t.Errorf("Error getting job by page ID: %v", err)
		}

		if res.ID != job.ID {
			t.Errorf("Expected ID to be %s, got %s", job.ID, res.ID)
		}

		if res.URL != job.URL {
			t.Errorf("Expected URL to be %s, got %s", job.URL, res.URL)
		}

		if res.CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, res.CreatedAt)
		}
	})

	t.Run("returns error if job does not exist", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		res, err := repo.GetByPageID("nonexistent")

		// Assert
		if err != domain.ErrCrawlJobNotFound {
			t.Errorf("Expected error, got %v", err)
		}

		if res != nil {
			t.Errorf("Expected nil, got %v", res)
		}
	})
}

func TestFirst(t *testing.T) {
	t.Run("can get first job", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			PageID:    "hash",
			CreatedAt: time.Now().UTC(),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.CreatedAt)

		job2 := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example2.com",
			PageID:    "hash2",
			CreatedAt: time.Now().UTC().Add(time.Second),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?)", job2.ID, job2.URL, job2.PageID, job2.CreatedAt)

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)
		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job2.ID)

		// Act
		res, err := repo.First()

		// Assert
		if err != nil {
			t.Errorf("Error getting first job: %v", err)
		}

		if res.ID != job.ID {
			t.Errorf("Expected ID to be %s, got %s", job.ID, res.ID)
		}

		if res.URL != job.URL {
			t.Errorf("Expected URL to be %s, got %s", job.URL, res.URL)
		}

		if res.PageID != job.PageID {
			t.Errorf("Expected PageID to be %s, got %s", job.PageID, res.PageID)
		}

		if res.CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, res.CreatedAt)
		}

	})

	t.Run("returns error if no jobs exist", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		res, err := repo.First()

		// Assert
		if err != domain.ErrCrawlJobNotFound {
			t.Errorf("Expected error, got %v", err)
		}

		if res != nil {
			t.Errorf("Expected nil, got %v", res)
		}
	})
}

func TestFirstProcessable(t *testing.T) {
	t.Run("can get first processable job", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			PageID:    "hash",
			CreatedAt: time.Now().UTC(),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.CreatedAt)

		job2 := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example2.com",
			PageID:    "hash2",
			CreatedAt: time.Now().UTC().Add(time.Second),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?)", job2.ID, job2.URL, job2.PageID, job2.CreatedAt)

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)
		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job2.ID)

		// Act
		res, err := repo.FirstProcessable()

		// Assert
		if err != nil {
			t.Errorf("Error getting first available job: %v", err)
		}

		if res.ID != job.ID {
			t.Errorf("Expected ID to be %s, got %s", job.ID, res.ID)
		}

		if res.URL != job.URL {
			t.Errorf("Expected URL to be %s, got %s", job.URL, res.URL)
		}

		if res.PageID != job.PageID {
			t.Errorf("Expected PageID to be %s, got %s", job.PageID, res.PageID)
		}

		if res.CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, res.CreatedAt)
		}

	})

	t.Run("does not return job where backoff_until is in the future", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:           util.UUID(),
			URL:          "http://example.com",
			PageID:       "hash",
			CreatedAt:    time.Now().UTC(),
			BackoffUntil: sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, backoff_until, created_at) VALUES (?, ?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.BackoffUntil, job.CreatedAt)

		job2 := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example2.com",
			PageID:    "hash2",
			CreatedAt: time.Now().UTC().Add(time.Second),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job2.ID, job2.URL, job2.PageID, job2.CreatedAt)

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)
		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job2.ID)

		// Act
		res, err := repo.FirstProcessable()

		// Assert
		if err != nil {
			t.Errorf("Error getting first available job: %v", err)
		}

		if res.ID != job2.ID {
			t.Errorf("Expected ID to be %s, got %s", job2.ID, res.ID)
		}

		if res.URL != job2.URL {
			t.Errorf("Expected URL to be %s, got %s", job2.URL, res.URL)
		}

		if res.PageID != job2.PageID {
			t.Errorf("Expected PageID to be %s, got %s", job2.PageID, res.PageID)
		}

		if res.CreatedAt.Sub(job2.CreatedAt) > time.Second || job2.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job2.CreatedAt, res.CreatedAt)
		}
	})

	t.Run("does not return job where last_crawled_at is within the last month", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:            "job1",
			URL:           "http://example.com",
			PageID:        "hash",
			CreatedAt:     time.Now().UTC(),
			LastCrawledAt: sql.NullTime{Time: time.Now().Add(-(time.Hour * 24 * 20)), Valid: true},
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, last_crawled_at, created_at) VALUES (?, ?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.LastCrawledAt, job.CreatedAt)

		job2 := &domain.CrawlJob{
			ID:        "job2",
			URL:       "http://example2.com",
			PageID:    "hash2",
			CreatedAt: time.Now().UTC().Add(time.Second),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job2.ID, job2.URL, job2.PageID, job2.CreatedAt)

		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)
		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job2.ID)

		// Act
		res, err := repo.FirstProcessable()

		// Assert
		if err != nil {
			t.Errorf("Error getting first job without backoff: %v", err)
		}

		if res.ID != job2.ID {
			t.Errorf("Expected ID to be %s, got %s", job2.ID, res.ID)
		}

		if res.URL != job2.URL {
			t.Errorf("Expected URL to be %s, got %s", job2.URL, res.URL)
		}

		if res.PageID != job2.PageID {
			t.Errorf("Expected PageID to be %s, got %s", job2.PageID, res.PageID)
		}

		if res.CreatedAt.Sub(job2.CreatedAt) > time.Second || job2.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, res.CreatedAt)
		}

	})

}

func TestDelete(t *testing.T) {
	t.Run("can delete a job", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			CreatedAt: time.Now().UTC(),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, page_id, created_at) VALUES (?, ?, ?, ?)", job.ID, job.URL, job.PageID, job.CreatedAt)
		defer db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)

		// Act
		err := repo.Delete(job.ID)

		// Assert
		if err != nil {
			t.Errorf("Error deleting job: %v", err)
		}

		res, err := db.Query("SELECT * FROM crawl_jobs WHERE id = ?", job.ID)

		if err != nil {
			t.Errorf("Error querying for job: %v", err)
		}

		var id string
		var url string
		var createdAt time.Time

		for res.Next() {
			err = res.Scan(&id, &url, &createdAt)
			if err != nil {
				t.Errorf("Error scanning job: %v", err)
			}
		}

		if id != "" {
			t.Errorf("Expected ID to be empty, got %s", id)
		}

		if url != "" {
			t.Errorf("Expected URL to be empty, got %s", url)
		}

		if !createdAt.IsZero() {
			t.Errorf("Expected CreatedAt to be zero, got %v", createdAt)
		}
	})
}

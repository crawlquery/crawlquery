package mysql_test

import (
	"crawlquery/api/crawl/job/repository/mysql"
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
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
			CreatedAt: time.Now().UTC(),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, created_at) VALUES (?, ?, ?)", job.ID, job.URL, job.CreatedAt)

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

		// Clean up
		db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)
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
		if err == nil {
			t.Errorf("Expected error, got nil")
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
			CreatedAt: time.Now().UTC(),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, created_at) VALUES (?, ?, ?)", job.ID, job.URL, job.CreatedAt)

		job2 := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example2.com",
			CreatedAt: time.Now().UTC().Add(time.Second),
		}
		db.Exec("INSERT INTO crawl_jobs (id, url, created_at) VALUES (?, ?, ?)", job2.ID, job2.URL, job2.CreatedAt)

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

		if res.CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, res.CreatedAt)
		}

		// Clean up
		db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job.ID)
		db.Exec("DELETE FROM crawl_jobs WHERE id = ?", job2.ID)
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
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if res != nil {
			t.Errorf("Expected nil, got %v", res)
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
		db.Exec("INSERT INTO crawl_jobs (id, url, created_at) VALUES (?, ?, ?)", job.ID, job.URL, job.CreatedAt)

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

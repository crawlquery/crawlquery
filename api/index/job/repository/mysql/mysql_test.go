package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/index/job/repository/mysql"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"database/sql"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create index job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			ID:     "job1-create-index-job",
			PageID: "page1",
			LastIndexedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			BackoffUntil: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			FailedReason: sql.NullString{
				String: "error",
				Valid:  true,
			},
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job.ID)

		_, err := repo.Create(job)

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		var checkJob domain.IndexJob

		err = db.QueryRow("SELECT id, page_id, backoff_until, last_indexed_at, failed_reason, created_at FROM index_jobs WHERE id = ?", job.ID).
			Scan(&checkJob.ID, &checkJob.PageID, &checkJob.BackoffUntil, &checkJob.LastIndexedAt, &checkJob.FailedReason, &checkJob.CreatedAt)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}

		if checkJob.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, checkJob.ID)
		}

		if checkJob.PageID != job.PageID {
			t.Errorf("Expected job PageID to be %s, got %s", job.PageID, checkJob.PageID)
		}

		if checkJob.BackoffUntil.Time.UTC().Round(time.Second) != job.BackoffUntil.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job BackoffUntil to be %s, got %s", job.BackoffUntil.Time, checkJob.BackoffUntil.Time)
		}

		if checkJob.LastIndexedAt.Time.UTC().Round(time.Second) != job.LastIndexedAt.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job LastIndexedAt to be %s, got %s", job.LastIndexedAt.Time, checkJob.LastIndexedAt.Time)
		}

		if checkJob.FailedReason.String != job.FailedReason.String {
			t.Errorf("Expected job FailedReason to be %s, got %s", job.FailedReason.String, checkJob.FailedReason.String)
		}

		if checkJob.CreatedAt.UTC().Round(time.Second) != job.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job CreatedAt to be %s, got %s", job.CreatedAt, checkJob.CreatedAt)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("can get index job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			ID:     "job1",
			PageID: "page1",
			LastIndexedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			BackoffUntil: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			FailedReason: sql.NullString{
				String: "error",
				Valid:  true,
			},
			CreatedAt: time.Now(),
		}

		db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			job.ID, job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.CreatedAt)

		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job.ID)

		result, err := repo.Get(job.ID)

		if err != nil {
			t.Fatalf("Error getting index job: %v", err)
		}

		if result.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, result.ID)
		}

		if result.PageID != job.PageID {
			t.Errorf("Expected job PageID to be %s, got %s", job.PageID, result.PageID)
		}

		if result.BackoffUntil.Time.UTC().Round(time.Second) != job.BackoffUntil.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job BackoffUntil to be %s, got %s", job.BackoffUntil.Time, result.BackoffUntil.Time)
		}

		if result.LastIndexedAt.Time.UTC().Round(time.Second) != job.LastIndexedAt.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job LastIndexedAt to be %s, got %s", job.LastIndexedAt.Time, result.LastIndexedAt.Time)
		}

		if result.FailedReason.String != job.FailedReason.String {
			t.Errorf("Expected job FailedReason to be %s, got %s", job.FailedReason.String, result.FailedReason.String)
		}

		if result.CreatedAt.UTC().Round(time.Second) != job.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job CreatedAt to be %s, got %s", job.CreatedAt, result.CreatedAt)
		}

	})
}

func TestGetByPageID(t *testing.T) {
	t.Run("can get index job by page ID", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			ID:     "job1-get-index-job-by-page-id",
			PageID: "page1",
			LastIndexedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			BackoffUntil: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			FailedReason: sql.NullString{
				String: "error",
				Valid:  true,
			},
			CreatedAt: time.Now(),
		}

		job2 := &domain.IndexJob{
			ID:        "job2-get-index-job-by-page-id",
			PageID:    "page2",
			CreatedAt: time.Now(),
		}

		_, err := db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			job.ID, job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting index job: %v", err)
		}

		_, err = db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			job2.ID, job2.PageID, job2.BackoffUntil, job2.LastIndexedAt, job2.FailedReason, job2.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting index job: %v", err)
		}

		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job.ID)
		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job2.ID)

		result, err := repo.GetByPageID(job.PageID)

		if err != nil {
			t.Fatalf("Error getting index job: %v", err)
		}

		if result.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, result.ID)
		}

		if result.PageID != job.PageID {
			t.Errorf("Expected job PageID to be %s, got %s", job.PageID, result.PageID)
		}

		if result.BackoffUntil.Time.UTC().Round(time.Second) != job.BackoffUntil.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job BackoffUntil to be %s, got %s", job.BackoffUntil.Time, result.BackoffUntil.Time)
		}

		if result.LastIndexedAt.Time.UTC().Round(time.Second) != job.LastIndexedAt.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job LastIndexedAt to be %s, got %s", job.LastIndexedAt.Time, result.LastIndexedAt.Time)
		}

		if result.FailedReason.String != job.FailedReason.String {
			t.Errorf("Expected job FailedReason to be %s, got %s", job.FailedReason.String, result.FailedReason.String)
		}

		if result.CreatedAt.UTC().Round(time.Second) != job.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job CreatedAt to be %s, got %s", job.CreatedAt, result.CreatedAt)
		}

	})
}

func TestUpdate(t *testing.T) {
	t.Run("can update index job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			ID:     "job1-can-update-index-job",
			PageID: "page1",
			LastIndexedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			BackoffUntil: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			FailedReason: sql.NullString{
				String: "error",
				Valid:  true,
			},
			CreatedAt: time.Now(),
		}

		db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			job.ID, job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.CreatedAt)

		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job.ID)

		job.PageID = "page2"
		job.BackoffUntil.Time = job.BackoffUntil.Time.Add(1 * time.Hour)
		job.LastIndexedAt.Time = job.LastIndexedAt.Time.Add(1 * time.Hour)
		job.FailedReason.String = "error2"

		err := repo.Update(job)

		if err != nil {
			t.Fatalf("Error updating index job: %v", err)
		}

		var checkJob domain.IndexJob

		err = db.QueryRow("SELECT id, page_id, backoff_until, last_indexed_at, failed_reason, created_at FROM index_jobs WHERE id = ?", job.ID).
			Scan(&checkJob.ID, &checkJob.PageID, &checkJob.BackoffUntil, &checkJob.LastIndexedAt, &checkJob.FailedReason, &checkJob.CreatedAt)

		if err != nil {
			t.Fatalf("Error getting index job: %v", err)
		}

		if checkJob.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, checkJob.ID)
		}

		if checkJob.PageID != job.PageID {
			t.Errorf("Expected job PageID to be %s, got %s", job.PageID, checkJob.PageID)
		}

		if checkJob.BackoffUntil.Time.UTC().Round(time.Second) != job.BackoffUntil.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job BackoffUntil to be %s, got %s", job.BackoffUntil.Time, checkJob.BackoffUntil.Time)
		}

		if checkJob.LastIndexedAt.Time.UTC().Round(time.Second) != job.LastIndexedAt.Time.UTC().Round(time.Second) {
			t.Errorf("Expected job LastIndexedAt to be %s, got %s", job.LastIndexedAt.Time, checkJob.LastIndexedAt.Time)
		}

		if checkJob.FailedReason.String != job.FailedReason.String {
			t.Errorf("Expected job FailedReason to be %s, got %s", job.FailedReason.String, checkJob.FailedReason.String)
		}

		if checkJob.CreatedAt.UTC().Round(time.Second) != job.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job CreatedAt to be %s, got %s", job.CreatedAt, checkJob.CreatedAt)
		}

	})
}

func TestNext(t *testing.T) {
	t.Run("can get next index job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			ID:        "job1-can-get-next-index-job",
			CreatedAt: time.Now(),
		}

		_, err := db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			job.ID, job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting index job: %v", err)
		}

		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job.ID)

		result, err := repo.Next()

		if err != nil {
			t.Fatalf("Error getting next index job: %v", err)
		}

		if result.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, result.ID)
		}
	})

	t.Run("does not return job if backoff is not expired", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			ID:        "job1-does-not-return-job",
			CreatedAt: time.Now(),
		}

		job.BackoffUntil.Valid = true
		job.BackoffUntil.Time = time.Now().Add(1 * time.Hour)

		_, err := db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			job.ID, job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting index job: %v", err)
		}

		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job.ID)

		result, err := repo.Next()

		if err != domain.ErrIndexJobNotFound {
			t.Errorf("Expected ErrIndexJobNotFound, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("does not return job if last indexed is set", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			ID:        "job1",
			CreatedAt: time.Now(),
		}

		job.LastIndexedAt.Time = time.Now()
		job.LastIndexedAt.Valid = true

		_, err := db.Exec("INSERT INTO index_jobs (id, page_id, backoff_until, last_indexed_at, failed_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			job.ID, job.PageID, job.BackoffUntil, job.LastIndexedAt, job.FailedReason, job.CreatedAt)

		if err != nil {
			t.Fatalf("Error inserting index job: %v", err)
		}

		defer db.Exec("DELETE FROM index_jobs WHERE id = ?", job.ID)

		result, err := repo.Next()

		if err != domain.ErrIndexJobNotFound {
			t.Errorf("Expected ErrIndexJobNotFound, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}

	})
}

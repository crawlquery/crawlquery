package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/index/job/repository/mysql"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	t.Run("can create index job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			PageID:    "page1",
			ShardID:   0,
			Status:    domain.IndexStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(job)
		defer db.Exec("DELETE FROM index_jobs WHERE page_id = ?", job.PageID)

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		var checkJob domain.IndexJob

		err = db.QueryRow("SELECT page_id, status, created_at, updated_at FROM index_jobs WHERE page_id = ?", job.PageID).Scan(&checkJob.PageID, &checkJob.Status, &checkJob.CreatedAt, &checkJob.UpdatedAt)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}
		if checkJob.PageID != job.PageID {
			t.Errorf("Expected job PageID to be %s, got %s", job.PageID, checkJob.PageID)
		}

		if checkJob.Status != job.Status {
			t.Errorf("Expected job Status to be %s, got %s", job.Status, checkJob.Status)
		}

		if checkJob.CreatedAt.UTC().Round(time.Second) != job.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job CreatedAt to be %s, got %s", job.CreatedAt, checkJob.CreatedAt)
		}

		if checkJob.UpdatedAt.UTC().Round(time.Second) != job.UpdatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job UpdatedAt to be %s, got %s", job.UpdatedAt, checkJob.UpdatedAt)
		}
	})

	t.Run("can update index job", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job := &domain.IndexJob{
			PageID:    "page1",
			ShardID:   0,
			Status:    domain.IndexStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(job)
		defer db.Exec("DELETE FROM index_jobs WHERE page_id = ?", job.PageID)

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		job.Status = domain.IndexStatusInProgress
		job.UpdatedAt = time.Now()

		err = repo.Save(job)

		if err != nil {
			t.Errorf("Error updating index job: %v", err)
		}

		var checkJob domain.IndexJob

		err = db.QueryRow("SELECT page_id, status, created_at, updated_at FROM index_jobs WHERE page_id = ?", job.PageID).Scan(&checkJob.PageID, &checkJob.Status, &checkJob.CreatedAt, &checkJob.UpdatedAt)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}
		if checkJob.PageID != job.PageID {
			t.Errorf("Expected job PageID to be %s, got %s", job.PageID, checkJob.PageID)
		}

		if checkJob.Status != job.Status {
			t.Errorf("Expected job Status to be %s, got %s", job.Status, checkJob.Status)
		}

		if checkJob.CreatedAt.UTC().Round(time.Second) != job.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job CreatedAt to be %s, got %s", job.CreatedAt, checkJob.CreatedAt)
		}

		if checkJob.UpdatedAt.UTC().Round(time.Second) != job.UpdatedAt.UTC().Round(time.Second) {
			t.Errorf("Expected job UpdatedAt to be %s, got %s", job.UpdatedAt, checkJob.UpdatedAt)
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
			PageID:    "page1",
			ShardID:   0,
			Status:    domain.IndexStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(job)
		defer db.Exec("DELETE FROM index_jobs WHERE page_id = ?", job.PageID)

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		result, err := repo.Get(job.PageID)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}

		if result.PageID != job.PageID {
			t.Errorf("Expected job ID to be %s, got %s", job.PageID, result.PageID)
		}
	})
}

func TestListByStatus(t *testing.T) {
	t.Run("can list index jobs by status", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		job1 := &domain.IndexJob{
			PageID:    "page1",
			Status:    domain.IndexStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		job2 := &domain.IndexJob{
			PageID:    "page2",
			Status:    domain.IndexStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		job3 := &domain.IndexJob{
			PageID:    "page3",
			Status:    domain.IndexStatusCompleted,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(job1)
		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}
		defer db.Exec("DELETE FROM index_jobs WHERE page_id = ?", job1.PageID)

		err = repo.Save(job2)
		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}
		defer db.Exec("DELETE FROM index_jobs WHERE page_id = ?", job2.PageID)

		err = repo.Save(job3)
		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}
		defer db.Exec("DELETE FROM index_jobs WHERE page_id = ?", job3.PageID)

		jobs, err := repo.ListByStatus(10, domain.IndexStatusPending)

		if err != nil {
			t.Errorf("Error listing index jobs: %v", err)
		}

		if len(jobs) != 2 {
			t.Errorf("Expected 2 jobs, got %v", len(jobs))
		}

		if jobs[0].PageID != job1.PageID {
			t.Errorf("Expected first job PageID to be %s, got %s", job1.PageID, jobs[0].PageID)
		}

		if jobs[1].PageID != job2.PageID {
			t.Errorf("Expected second job PageID to be %s, got %s", job2.PageID, jobs[1].PageID)
		}

		jobs, err = repo.ListByStatus(10, domain.IndexStatusCompleted)

		if err != nil {
			t.Errorf("Error listing index jobs: %v", err)
		}

		if len(jobs) != 1 {
			t.Errorf("Expected 1 job, got %v", len(jobs))
		}

		if jobs[0].PageID != job3.PageID {
			t.Errorf("Expected job PageID to be %s, got %s", job3.PageID, jobs[0].PageID)
		}

	})
}

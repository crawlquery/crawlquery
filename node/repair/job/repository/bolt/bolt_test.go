package bolt_test

import (
	repairJobRepo "crawlquery/node/repair/job/repository/bolt"
	"time"

	"crawlquery/node/domain"
	"testing"

	"github.com/boltdb/bolt"
)

func TestCreate(t *testing.T) {
	t.Run("can create job", func(t *testing.T) {

		db, err := bolt.Open("/tmp/repair_job_test.db", 0600, nil)
		if err != nil {
			t.Fatalf("Error opening db: %v", err)
		}
		defer db.Close()

		repo := repairJobRepo.NewRepository(db)
		job := &domain.RepairJob{
			PageID: "1",
		}

		err = repo.Create(job)

		if err != nil {
			t.Fatalf("Error creating job: %v", err)
		}

		check, err := repo.Get("1")

		if err != nil {
			t.Fatalf("Error getting job: %v", err)
		}

		if check.PageID != "1" {
			t.Fatalf("Expected job ID to be '1', got '%s'", check.PageID)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("can update job", func(t *testing.T) {

		db, err := bolt.Open("/tmp/repair_job_test.db", 0600, nil)
		if err != nil {
			t.Fatalf("Error opening db: %v", err)
		}
		defer db.Close()

		repo := repairJobRepo.NewRepository(db)
		job := &domain.RepairJob{
			PageID: "1",
		}

		err = repo.Create(job)

		if err != nil {
			t.Fatalf("Error creating job: %v", err)
		}

		job.Status = domain.RepairJobStatusComplete
		job.StatusLastUpdatedAt = time.Now()

		err = repo.Update(job)

		if err != nil {
			t.Fatalf("Error updating job: %v", err)
		}

		check, err := repo.Get("1")

		if err != nil {
			t.Fatalf("Error getting job: %v", err)
		}

		if check.PageID != "1" {
			t.Fatalf("Expected job ID to be '1', got '%s'", check.PageID)
		}

		if check.Status != domain.RepairJobStatusComplete {
			t.Fatalf("Expected job status to be 'Complete', got '%s'", check.Status)
		}

		expectedTime := job.StatusLastUpdatedAt.Round(time.Second)

		if check.StatusLastUpdatedAt.Round(time.Second) != expectedTime {
			t.Fatalf("Expected job status last updated at to be '%v', got '%v'", expectedTime, check.StatusLastUpdatedAt)
		}
	})
}

package mem

import (
	"crawlquery/node/domain"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create job", func(t *testing.T) {
		repo := NewRepository()
		job := &domain.RepairJob{
			PageID: "1",
		}

		err := repo.Create(job)

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
		repo := NewRepository()
		job := &domain.RepairJob{
			PageID: "1",
		}

		err := repo.Create(job)

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

		if check.Status != domain.RepairJobStatusComplete {
			t.Fatalf("Expected job status to be 'complete', got '%s'", check.Status)
		}

		if check.StatusLastUpdatedAt != job.StatusLastUpdatedAt {
			t.Fatalf("Expected job status last updated at to be '%v', got '%v'", job.StatusLastUpdatedAt, check.StatusLastUpdatedAt)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("can get job", func(t *testing.T) {
		repo := NewRepository()
		job := &domain.RepairJob{
			PageID: "1",
		}

		err := repo.Create(job)

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

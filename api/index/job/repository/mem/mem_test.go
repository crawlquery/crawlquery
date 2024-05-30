package mem

import (
	"crawlquery/api/domain"
	"testing"
)

func TestSave(t *testing.T) {
	t.Run("can save index job", func(t *testing.T) {
		repo := NewRepository()

		job := &domain.IndexJob{
			PageID: "job1",
		}

		err := repo.Save(job)

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		if len(repo.jobs) != 1 {
			t.Errorf("Expected 1 job, got %v", len(repo.jobs))
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("can get index job", func(t *testing.T) {
		repo := NewRepository()

		job := &domain.IndexJob{
			PageID: "job1",
		}

		repo.jobs[job.PageID] = job

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
		repo := NewRepository()

		job1 := &domain.IndexJob{
			PageID: "job1",
			Status: domain.IndexStatusPending,
		}

		job2 := &domain.IndexJob{
			PageID: "job2",
			Status: domain.IndexStatusPending,
		}

		job3 := &domain.IndexJob{
			PageID: "job3",
			Status: domain.IndexStatusCompleted,
		}

		repo.jobs[job1.PageID] = job1
		repo.jobs[job2.PageID] = job2
		repo.jobs[job3.PageID] = job3

		jobs, err := repo.ListByStatus(10, domain.IndexStatusPending)

		if err != nil {
			t.Errorf("Error listing index jobs: %v", err)
		}

		if len(jobs) != 2 {
			t.Errorf("Expected 2 jobs, got %v", len(jobs))
		}
	})
}

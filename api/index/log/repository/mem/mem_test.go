package mem

import (
	"crawlquery/api/domain"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	t.Run("can save index job", func(t *testing.T) {
		repo := NewRepository()

		log := &domain.IndexLog{
			PageID:    "job1",
			Status:    domain.IndexStatusPending,
			Info:      "info",
			CreatedAt: time.Now(),
		}

		err := repo.Save(log)

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		if len(repo.logs) != 1 {
			t.Errorf("Expected 1 job, got %v", len(repo.logs))
		}
	})
}

func TestListByPageID(t *testing.T) {
	t.Run("can list index jobs by pageID", func(t *testing.T) {
		repo := NewRepository()

		log1 := &domain.IndexLog{
			ID:     "log1",
			PageID: "page1",
			Status: domain.IndexStatusPending,
			Info:   "info",
		}

		log2 := &domain.IndexLog{
			ID:     "log2",
			PageID: "page2",
			Status: domain.IndexStatusPending,
			Info:   "info",
		}

		log3 := &domain.IndexLog{
			ID:     "log3",
			PageID: "page3",
			Status: domain.IndexStatusCompleted,
			Info:   "info",
		}

		err := repo.Save(log1)
		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		err = repo.Save(log2)
		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		err = repo.Save(log3)
		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		logs, err := repo.ListByPageID("page1")

		if err != nil {
			t.Errorf("Error listing index jobs: %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("Expected 1 jobs, got %v", len(logs))
		}
	})
}

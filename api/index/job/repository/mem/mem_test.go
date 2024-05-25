package mem

import (
	"crawlquery/api/domain"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create index job", func(t *testing.T) {
		repo := NewRepository()

		job := &domain.IndexJob{
			ID: "job1",
		}

		_, err := repo.Create(job)

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
			ID: "job1",
		}

		repo.jobs[job.ID] = job

		result, err := repo.Get(job.ID)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}

		if result.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, result.ID)
		}
	})
}

func TestNext(t *testing.T) {
	t.Run("can get next index job", func(t *testing.T) {
		repo := NewRepository()

		job := &domain.IndexJob{
			ID: "job1",
		}

		repo.jobs[job.ID] = job

		result, err := repo.Next()

		if err != nil {
			t.Fatalf("Error getting next index job: %v", err)
		}

		if result.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, result.ID)
		}
	})

	t.Run("does not return job if backoff is not expired", func(t *testing.T) {
		repo := NewRepository()

		job := &domain.IndexJob{
			ID: "job1",
		}

		job.BackoffUntil.Valid = true
		job.BackoffUntil.Time = job.BackoffUntil.Time.Add(1 * time.Hour)

		repo.jobs[job.ID] = job

		result, err := repo.Next()

		if err != domain.ErrIndexJobNotFound {
			t.Errorf("Expected ErrIndexJobNotFound, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("does not return job if last indexed is set", func(t *testing.T) {
		repo := NewRepository()

		job := &domain.IndexJob{
			ID: "job1",
		}

		job.LastIndexedAt.Time = time.Now()
		job.LastIndexedAt.Valid = true

		repo.jobs[job.ID] = job

		result, err := repo.Next()

		if err != domain.ErrIndexJobNotFound {
			t.Errorf("Expected ErrIndexJobNotFound, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("can update index job", func(t *testing.T) {
		repo := NewRepository()

		job := &domain.IndexJob{
			ID: "job1",
		}

		repo.jobs[job.ID] = job

		job.LastIndexedAt.Time = time.Now()
		job.LastIndexedAt.Valid = true

		err := repo.Update(job)

		if err != nil {
			t.Errorf("Error updating index job: %v", err)
		}

		if !repo.jobs[job.ID].LastIndexedAt.Valid {
			t.Errorf("Expected LastIndexedAt to be valid")
		}
	})
}

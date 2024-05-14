package mem

import (
	"crawlquery/api/domain"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	repo := NewRepository()

	job := &domain.CrawlJob{
		ID:  "job1",
		URL: "http://example.com",
	}

	err := repo.Create(job)

	if err != nil {
		t.Fatalf("Error creating job: %v", err)
	}

	if repo.jobs[job.ID].ID != job.ID {
		t.Errorf("Expected ID to be %s, got %s", job.ID, repo.jobs[job.ID].ID)
	}

	if repo.jobs[job.ID].URL != job.URL {
		t.Errorf("Expected URL to be %s, got %s", job.URL, repo.jobs[job.ID].URL)
	}

	if repo.jobs[job.ID].CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(repo.jobs[job.ID].CreatedAt) > time.Second {
		t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, repo.jobs[job.ID].CreatedAt)
	}
}

func TestUpdate(t *testing.T) {
	repo := NewRepository()

	job := &domain.CrawlJob{
		ID:  "job1",
		URL: "http://example.com",
	}

	repo.jobs[job.ID] = job

	job.URL = "http://example2.com"

	err := repo.Update(job)

	if err != nil {
		t.Fatalf("Error updating job: %v", err)
	}

	if repo.jobs[job.ID].ID != job.ID {
		t.Errorf("Expected ID to be %s, got %s", job.ID, repo.jobs[job.ID].ID)
	}

	if repo.jobs[job.ID].URL != job.URL {
		t.Errorf("Expected URL to be %s, got %s", job.URL, repo.jobs[job.ID].URL)
	}

	if repo.jobs[job.ID].CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(repo.jobs[job.ID].CreatedAt) > time.Second {
		t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, repo.jobs[job.ID].CreatedAt)
	}
}

func TestGet(t *testing.T) {
	repo := NewRepository()

	job := &domain.CrawlJob{
		ID:  "job1",
		URL: "http://example.com",
	}

	repo.jobs[job.ID] = job

	got, err := repo.Get(job.ID)

	if err != nil {
		t.Fatalf("Error getting job: %v", err)
	}

	if got.ID != job.ID {
		t.Errorf("Expected ID to be %s, got %s", job.ID, got.ID)
	}

	if got.URL != job.URL {
		t.Errorf("Expected URL to be %s, got %s", job.URL, got.URL)
	}

	if got.CreatedAt.Sub(job.CreatedAt) > time.Second || job.CreatedAt.Sub(got.CreatedAt) > time.Second {
		t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job.CreatedAt, got.CreatedAt)
	}
}

func TestFirst(t *testing.T) {
	repo := NewRepository()

	job1 := &domain.CrawlJob{
		ID:  "job1",
		URL: "http://example.com",
	}

	job2 := &domain.CrawlJob{
		ID:  "job2",
		URL: "http://example.com",
	}

	repo.jobs[job1.ID] = job1
	repo.jobs[job2.ID] = job2

	first, err := repo.First()

	if err != nil {
		t.Fatalf("Error getting first job: %v", err)
	}

	if first.ID != job1.ID {
		t.Errorf("Expected ID to be %s, got %s", job1.ID, first.ID)
	}

	if first.URL != job1.URL {
		t.Errorf("Expected URL to be %s, got %s", job1.URL, first.URL)
	}

	if first.CreatedAt.Sub(job1.CreatedAt) > time.Second || job1.CreatedAt.Sub(first.CreatedAt) > time.Second {
		t.Errorf("Expected CreatedAt to be within one second of %v, got %v", job1.CreatedAt, first.CreatedAt)
	}
}

func TestDelete(t *testing.T) {
	repo := NewRepository()

	job := &domain.CrawlJob{
		ID:  "job1",
		URL: "http://example.com",
	}

	repo.jobs[job.ID] = job

	err := repo.Delete(job.ID)

	if err != nil {
		t.Fatalf("Error deleting job: %v", err)
	}

	if _, ok := repo.jobs[job.ID]; ok {
		t.Errorf("Expected job to be deleted, got %v", repo.jobs[job.ID])
	}
}

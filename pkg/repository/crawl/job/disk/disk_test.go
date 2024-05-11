package disk_test

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/repository/crawl/job/disk"
	"os"
	"testing"
)

func TestDisk(t *testing.T) {
	defer os.Remove("/tmp/crawl_job.gob")
	filepath := "/tmp/crawl_job.gob"
	repo := disk.NewDiskRepository(filepath)

	err := repo.Push(&domain.CrawlJob{
		URL: "http://google.com",
	})

	if err != nil {
		t.Fatalf("Error pushing to disk repository: %v", err)
	}

	err = repo.Push(&domain.CrawlJob{
		URL: "http://facebook.com",
	})

	if err != nil {
		t.Fatalf("Error pushing to disk repository: %v", err)
	}

	job, err := repo.Pop()

	if err != nil {
		t.Fatalf("Error popping from disk repository: %v", err)
	}

	if job.URL != "http://google.com" {
		t.Errorf("Expected URL to be http://google.com, got %v", job.URL)
	}

	repoB := disk.NewDiskRepository(filepath)

	repoB.Load()

	job, err = repoB.Pop()

	if err != nil {
		t.Fatalf("Error popping from disk repository: %v", err)
	}

	if job.URL != "http://facebook.com" {
		t.Errorf("Expected URL to be http://facebook.com, got %v", job.URL)
	}

	job, err = repoB.Pop()

	if err == nil {
		t.Fatalf("Expected error popping from disk repository, got nil")
	}

	if job != nil {
		t.Errorf("Expected job to be nil, got %v", job)
	}
}

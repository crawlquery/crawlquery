package mem

import (
	"crawlquery/api/domain"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("can get a job", func(t *testing.T) {
		crawlJobRepo := NewRepository()
		crawlJob := &domain.CrawlJob{
			PageID: "page1",
			Status: domain.CrawlStatusPending,
		}

		crawlJobRepo.jobs["page1"] = crawlJob

		job, err := crawlJobRepo.Get("page1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if job.PageID != "page1" {
			t.Errorf("expected page1, got %v", job.PageID)
		}

		if job.Status != domain.CrawlStatusPending {
			t.Errorf("expected pending, got %v", job.Status)
		}
	})

	t.Run("returns error if job not found", func(t *testing.T) {
		crawlJobRepo := NewRepository()

		_, err := crawlJobRepo.Get("page1")

		if err != domain.ErrCrawlJobNotFound {
			t.Errorf("expected ErrCrawlJobNotFound, got %v", err)
		}
	})

}

func TestSave(t *testing.T) {
	t.Run("can save a job", func(t *testing.T) {
		crawlJobRepo := NewRepository()
		crawlJob := &domain.CrawlJob{
			PageID: "page1",
			Status: domain.CrawlStatusPending,
		}

		err := crawlJobRepo.Save(crawlJob)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if crawlJobRepo.jobs["page1"] != crawlJob {
			t.Errorf("expected job to be saved")
		}
	})
}

func TestListByStatus(t *testing.T) {
	t.Run("can list jobs by status", func(t *testing.T) {
		crawlJobRepo := NewRepository()
		for i := 0; i < 5; i++ {
			crawlJob := &domain.CrawlJob{
				PageID: domain.PageID(fmt.Sprintf("page%d", i)),
				Status: domain.CrawlStatusPending,
			}
			crawlJobRepo.jobs[crawlJob.PageID] = crawlJob
		}

		for i := 5; i < 10; i++ {
			crawlJob := &domain.CrawlJob{
				PageID: domain.PageID(fmt.Sprintf("page%d", i+5)),
				Status: domain.CrawlStatusInProgress,
			}
			crawlJobRepo.jobs[crawlJob.PageID] = crawlJob
		}

		jobs, err := crawlJobRepo.ListByStatus(3, domain.CrawlStatusPending)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(jobs) != 3 {
			t.Errorf("expected 3 jobs, got %v", len(jobs))
		}

		for _, job := range jobs {
			if job.Status != domain.CrawlStatusPending {
				t.Errorf("expected pending, got %v", job.Status)
			}
		}
	})
}

package mem

import (
	"crawlquery/api/domain"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("can get a job", func(t *testing.T) {
		crawlJobRepo := NewRepository()
		crawlJob := &domain.CrawlJob{
			PageID: "page1",
			Status: domain.CrawlJobStatusPending,
		}

		crawlJobRepo.jobs["page1"] = crawlJob

		job, err := crawlJobRepo.Get("page1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if job.PageID != "page1" {
			t.Errorf("expected page1, got %v", job.PageID)
		}

		if job.Status != domain.CrawlJobStatusPending {
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
			Status: domain.CrawlJobStatusPending,
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

package service_test

import (
	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlLogRepo "crawlquery/api/crawl/log/repository/mem"
	crawlService "crawlquery/api/crawl/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreateJob(t *testing.T) {
	t.Run("should create a job", func(t *testing.T) {
		crawlRepo := crawlJobRepo.NewRepository()
		crawlLogRepo := crawlLogRepo.NewRepository()
		crawlService := crawlService.NewService(
			crawlService.WithCrawlJobRepo(crawlRepo),
			crawlService.WithCrawlLogRepo(
				crawlLogRepo,
			),
			crawlService.WithLogger(testutil.NewTestLogger()),
		)

		pageID := util.PageID("http://example.com")

		err := crawlService.CreateJob(pageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		createdJob, err := crawlRepo.Get(pageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if createdJob.Status != domain.CrawlStatusPending {
			t.Errorf("expected status to be pending, got %v", createdJob.Status)
		}

		if createdJob.PageID != pageID {
			t.Errorf("expected pageID to be %v, got %v", pageID, createdJob.PageID)
		}

		logs, err := crawlLogRepo.ListByPageID(pageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("expected 1 log, got %v", len(logs))
		}

		if logs[0].PageID != pageID {
			t.Errorf("expected log pageID to be %v, got %v", pageID, logs[0].PageID)
		}

		if logs[0].Status != domain.CrawlStatusPending {
			t.Errorf("expected log status to be pending, got %v", logs[0].Status)
		}
	})
}

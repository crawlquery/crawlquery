package service_test

import (
	"context"
	"crawlquery/api/testfactory"
	"crawlquery/node/dto"
	"time"

	"crawlquery/api/domain"

	"crawlquery/pkg/util"
	"testing"

	"github.com/h2non/gock"
)

func TestCreateJob(t *testing.T) {
	t.Run("should create a job", func(t *testing.T) {
		sf := testfactory.NewServiceFactory()

		page := &domain.Page{
			ID:  util.PageID("http://example.com"),
			URL: "http://example.com",
		}

		err := sf.CrawlService.CreateJob(page)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		createdJob, err := sf.CrawlJobRepo.Get(page.ID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if createdJob.Status != domain.CrawlStatusPending {
			t.Errorf("expected status to be pending, got %v", createdJob.Status)
		}

		if createdJob.PageID != page.ID {
			t.Errorf("expected PageID to be %v, got %v", page.ID, createdJob.PageID)
		}

		logs, err := sf.CrawlLogRepo.ListByPageID(page.ID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("expected 1 log, got %v", len(logs))
		}

		if logs[0].PageID != page.ID {
			t.Errorf("expected log pageID to be %v, got %v", page.ID, logs[0].PageID)
		}

		if logs[0].Status != domain.CrawlStatusPending {
			t.Errorf("expected log status to be pending, got %v", logs[0].Status)
		}
	})
}

func TestProcessQueueItem(t *testing.T) {
	t.Run("should send crawl request to node", func(t *testing.T) {

		job := &domain.CrawlJob{
			PageID:  util.PageID("http://example.com"),
			URL:     "http://example.com",
			ShardID: 0,
			Status:  domain.CrawlStatusPending,
		}
		sf := testfactory.NewServiceFactory(
			testfactory.WithShard(&domain.Shard{ID: 0}),
			testfactory.WithNode(&domain.Node{
				ID:        "node1",
				ShardID:   0,
				Hostname:  "node1.cluster.com",
				Port:      8080,
				CreatedAt: time.Now(),
			}),
			testfactory.WithCrawlJob(job),
		)

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			Reply(200).
			JSON(dto.CrawlResponse{
				Links: []string{
					"http://example.com/1",
					"http://example.com/2",
				},
			})

		ctx := context.Background()
		sf.CrawlService.CacheNodes()
		err := sf.CrawlService.ProcessQueueItem(ctx, job)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("should update job status to completed", func(t *testing.T) {
		job := &domain.CrawlJob{
			PageID:  util.PageID("http://example.com"),
			URL:     "http://example.com",
			ShardID: 0,
			Status:  domain.CrawlStatusPending,
		}

		sf := testfactory.NewServiceFactory(
			testfactory.WithNode(&domain.Node{
				ID:        "node1",
				ShardID:   0,
				Hostname:  "node1.cluster.com",
				Port:      8080,
				CreatedAt: time.Now(),
			}),
			testfactory.WithCrawlJob(job),
		)

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			Reply(200).
			JSON(dto.CrawlResponse{
				Links: []string{
					"http://example.com/1",
					"http://example.com/2",
				},
			})

		ctx := context.Background()
		sf.CrawlService.CacheNodes()
		err := sf.CrawlService.ProcessQueueItem(ctx, job)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		updatedJob, err := sf.CrawlJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if updatedJob.Status != domain.CrawlStatusCompleted {
			t.Errorf("expected job status to be completed, got %v", updatedJob.Status)
		}
	})

	t.Run("should update job status to failed", func(t *testing.T) {

		job := &domain.CrawlJob{
			PageID:  util.PageID("http://example.com"),
			URL:     "http://example.com",
			ShardID: 0,
			Status:  domain.CrawlStatusPending,
		}

		sf := testfactory.NewServiceFactory(
			testfactory.WithNode(&domain.Node{
				ID:        "node1",
				ShardID:   0,
				Hostname:  "node1.cluster.com",
				Port:      8080,
				CreatedAt: time.Now(),
			}),
			testfactory.WithCrawlJob(job),
		)

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			Reply(500)

		ctx := context.Background()
		sf.CrawlService.CacheNodes()
		err := sf.CrawlService.ProcessQueueItem(ctx, job)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		updatedJob, err := sf.CrawlJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if updatedJob.Status != domain.CrawlStatusFailed {
			t.Errorf("expected job status to be failed, got %v", updatedJob.Status)
		}
	})
}

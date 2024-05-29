package service_test

import (
	"context"
	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlLogRepo "crawlquery/api/crawl/log/repository/mem"
	crawlService "crawlquery/api/crawl/service"
	"crawlquery/node/dto"
	"time"

	"crawlquery/api/domain"
	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"

	"github.com/h2non/gock"
)

func setup() map[string]interface{} {

	nodeRepo := nodeRepo.NewRepository()
	nodeService := nodeService.NewService(
		nodeService.WithNodeRepo(nodeRepo),
		nodeService.WithLogger(testutil.NewTestLogger()),
	)

	crawlRepo := crawlJobRepo.NewRepository()
	crawlLogRepo := crawlLogRepo.NewRepository()
	crawlService := crawlService.NewService(
		crawlService.WithCrawlJobRepo(crawlRepo),
		crawlService.WithNodeService(nodeService),
		crawlService.WithCrawlLogRepo(
			crawlLogRepo,
		),
		crawlService.WithLogger(testutil.NewTestLogger()),
	)

	return map[string]interface{}{
		"nodeRepo":     nodeRepo,
		"nodeService":  nodeService,
		"crawlRepo":    crawlRepo,
		"crawlLogRepo": crawlLogRepo,
		"crawlService": crawlService,
	}
}

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

		page := &domain.Page{
			ID:  util.PageID("http://example.com"),
			URL: "http://example.com",
		}

		err := crawlService.CreateJob(page)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		createdJob, err := crawlRepo.Get(page.ID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if createdJob.Status != domain.CrawlStatusPending {
			t.Errorf("expected status to be pending, got %v", createdJob.Status)
		}

		if createdJob.PageID != page.ID {
			t.Errorf("expected PageID to be %v, got %v", page.ID, createdJob.PageID)
		}

		logs, err := crawlLogRepo.ListByPageID(page.ID)

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
		ifs := setup()

		nodeRepo := ifs["nodeRepo"].(*nodeRepo.Repository)

		nodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		crawlJobRepo := ifs["crawlRepo"].(*crawlJobRepo.Repository)

		job := &domain.CrawlJob{
			PageID:  util.PageID("http://example.com"),
			URL:     "http://example.com",
			ShardID: 0,
			Status:  domain.CrawlStatusPending,
		}

		crawlJobRepo.Save(job)

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

		crawlService := ifs["crawlService"].(*crawlService.Service)

		ctx := context.Background()
		crawlService.CacheNodes()
		err := crawlService.ProcessQueueItem(ctx, job)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("should update job status to completed", func(t *testing.T) {
		ifs := setup()

		nodeRepo := ifs["nodeRepo"].(*nodeRepo.Repository)

		nodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		crawlJobRepo := ifs["crawlRepo"].(*crawlJobRepo.Repository)

		job := &domain.CrawlJob{
			PageID:  util.PageID("http://example.com"),
			URL:     "http://example.com",
			ShardID: 0,
			Status:  domain.CrawlStatusPending,
		}

		crawlJobRepo.Save(job)

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

		crawlService := ifs["crawlService"].(*crawlService.Service)

		ctx := context.Background()
		crawlService.CacheNodes()
		err := crawlService.ProcessQueueItem(ctx, job)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		updatedJob, err := crawlJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if updatedJob.Status != domain.CrawlStatusCompleted {
			t.Errorf("expected job status to be completed, got %v", updatedJob.Status)
		}
	})

	t.Run("should update job status to failed", func(t *testing.T) {
		ifs := setup()

		nodeRepo := ifs["nodeRepo"].(*nodeRepo.Repository)

		nodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		crawlJobRepo := ifs["crawlRepo"].(*crawlJobRepo.Repository)

		job := &domain.CrawlJob{
			PageID:  util.PageID("http://example.com"),
			URL:     "http://example.com",
			ShardID: 0,
			Status:  domain.CrawlStatusPending,
		}

		crawlJobRepo.Save(job)

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			Reply(500)

		crawlService := ifs["crawlService"].(*crawlService.Service)

		ctx := context.Background()
		crawlService.CacheNodes()
		err := crawlService.ProcessQueueItem(ctx, job)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		updatedJob, err := crawlJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if updatedJob.Status != domain.CrawlStatusFailed {
			t.Errorf("expected job status to be failed, got %v", updatedJob.Status)
		}
	})
}

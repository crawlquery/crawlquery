package service_test

import (
	"context"
	"crawlquery/api/testfactory"
	"crawlquery/node/dto"
	"fmt"
	"strings"
	"time"

	"crawlquery/api/domain"

	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"

	"github.com/h2non/gock"

	crawlService "crawlquery/api/crawl/service"
	crawlThrottleService "crawlquery/api/crawl/throttle/service"
)

func setupCrawlTests() (*testfactory.ServiceFactory, *domain.CrawlJob, *domain.Node) {
	job := &domain.CrawlJob{
		PageID:  util.PageID("http://example.com"),
		URL:     "http://example.com",
		ShardID: 0,
		Status:  domain.CrawlStatusPending,
	}

	node := &domain.Node{
		ID:        "node1",
		ShardID:   0,
		Hostname:  "node1.cluster.com",
		Port:      8080,
		CreatedAt: time.Now(),
	}

	sf := testfactory.NewServiceFactory(
		testfactory.WithShard(&domain.Shard{ID: 0}),
		testfactory.WithCrawlJob(job),
	)

	return sf, job, node
}

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
		sf, job, node := setupCrawlTests()

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
		err := sf.CrawlService.ProcessQueueItem(ctx, job, node)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("should update job status to completed", func(t *testing.T) {
		sf, job, node := setupCrawlTests()

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
		err := sf.CrawlService.ProcessQueueItem(ctx, job, node)

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

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("handles 400 error from node", func(t *testing.T) {
		sf, job, node := setupCrawlTests()

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			Reply(400).
			JSON(dto.ErrorResponse{
				Error: "request timeout error",
			})

		ctx := context.Background()
		err := sf.CrawlService.ProcessQueueItem(ctx, job, node)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "unexpected status code: 400 (request timeout error") {
			t.Errorf("expected error, got %v", err)
		}

		updatedJob, err := sf.CrawlJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if updatedJob.Status != domain.CrawlStatusFailed {
			t.Errorf("expected job status to be failed, got %v", updatedJob.Status)
		}

		logs, err := sf.CrawlLogRepo.ListByPageID(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("expected 1 logs, got %v", len(logs))
		}

		if logs[0].PageID != job.PageID {
			t.Errorf("expected log pageID to be %v, got %v", job.PageID, logs[0].PageID)
		}

		if logs[0].Status != domain.CrawlStatusFailed {
			t.Errorf("expected log status to be failed, got %v", logs[0].Status)
		}

		if logs[0].Info != "unexpected status code: 400 (request timeout error)" {
			t.Errorf("expected log info to be 'request timeout error', got %v", logs[0].Info)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("handles 500 error from node", func(t *testing.T) {
		sf, job, node := setupCrawlTests()

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			Reply(500)

		ctx := context.Background()
		err := sf.CrawlService.ProcessQueueItem(ctx, job, node)

		if !strings.Contains("unexpected status code: 500", err.Error()) {
			t.Errorf("expected error, got nil")
		}

		updatedJob, err := sf.CrawlJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if updatedJob.Status != domain.CrawlStatusFailed {
			t.Errorf("expected job status to be failed, got %v", updatedJob.Status)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("should timeout after deadline", func(t *testing.T) {
		sf, job, node := setupCrawlTests()

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			Reply(200).
			JSON(dto.CrawlResponse{
				Links: []string{
					"http://example.com/1",
					"http://example.com/2",
				},
			}).
			Delay(time.Second * 20)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		err := sf.CrawlService.ProcessQueueItem(ctx, job, node)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
			t.Errorf("expected error to be context.DeadlineExceeded, got %v", err)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})
}

func TestRunCrawlProcess(t *testing.T) {
	t.Run("should process crawl jobs with 10 workers and 100 crawl jobs", func(t *testing.T) {
		sf := testfactory.NewServiceFactory()

		crawlThrottleService := crawlThrottleService.NewService(
			crawlThrottleService.WithRateLimit(time.Second * 20),
		)
		crawlService := crawlService.NewService(
			crawlService.WithEventService(sf.EventService),
			crawlService.WithCrawlJobRepo(sf.CrawlJobRepo),
			crawlService.WithNodeService(sf.NodeService),
			crawlService.WithCrawlThrottleService(
				crawlThrottleService,
			),
			crawlService.WithCrawlLogRepo(
				sf.CrawlLogRepo,
			),
			crawlService.WithLogger(testutil.NewTestLogger()),
			crawlService.WithWorkers(10),
			crawlService.WithMaxQueueSize(100),
		)

		sf.ShardRepo.Create(&domain.Shard{ID: 0})
		sf.ShardRepo.Create(&domain.Shard{ID: 1})

		sf.NodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "shard0.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		sf.NodeRepo.Create(&domain.Node{
			ID:        "node2",
			ShardID:   1,
			Hostname:  "shard1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		shardJobs := make(map[domain.ShardID][]*domain.CrawlJob)

		defer gock.Off()

		for i := 0; i < 100; i++ {
			url := domain.URL(fmt.Sprintf("http://example%d.com/about", i))
			shardID, err := sf.ShardService.GetURLShardID(url)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			job := &domain.CrawlJob{
				PageID:  util.PageID(url),
				URL:     url,
				ShardID: shardID,
				Status:  domain.CrawlStatusPending,
			}

			err = sf.CrawlJobRepo.Save(job)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			shardJobs[shardID] = append(shardJobs[shardID], job)

			hostname := fmt.Sprintf("shard%d.cluster.com", shardID)

			status := 200

			if i == 50 {
				status = 500
			}
			gock.New(hostname).
				Post("/crawl").
				JSON(dto.CrawlRequest{
					PageID: string(job.PageID),
					URL:    string(job.URL),
				}).
				Reply(status).
				JSON(dto.CrawlResponse{
					Links: []string{
						"http://example.com/1",
						"http://example.com/2",
					},
				})
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		err := crawlService.RunCrawlProcess(ctx)

		if err != context.DeadlineExceeded {
			t.Errorf("expected got context deadline exceeded %v", err)
		}

		for _, jobs := range shardJobs {
			for _, job := range jobs {
				updatedJob, err := sf.CrawlJobRepo.Get(job.PageID)

				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if job.URL == "http://example50.com/about" {
					if updatedJob.Status != domain.CrawlStatusFailed {
						t.Errorf("expected job status to be failed, got %v", updatedJob.Status)
					}
				} else {
					if updatedJob.Status != domain.CrawlStatusCompleted {
						t.Errorf("expected job status to be completed, got %v", updatedJob.Status)
					}
				}
			}
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("should throttle urls of the same domain", func(t *testing.T) {
		sf := testfactory.NewServiceFactory()

		crawlThrottleService := crawlThrottleService.NewService(
			crawlThrottleService.WithRateLimit(time.Second * 20),
		)
		crawlService := crawlService.NewService(
			crawlService.WithEventService(sf.EventService),
			crawlService.WithCrawlJobRepo(sf.CrawlJobRepo),
			crawlService.WithNodeService(sf.NodeService),
			crawlService.WithCrawlThrottleService(
				crawlThrottleService,
			),
			crawlService.WithCrawlLogRepo(
				sf.CrawlLogRepo,
			),
			crawlService.WithLogger(testutil.NewTestLogger()),
			crawlService.WithWorkers(10),
			crawlService.WithMaxQueueSize(100),
		)

		sf.ShardRepo.Create(&domain.Shard{ID: 0})
		sf.ShardRepo.Create(&domain.Shard{ID: 1})

		sf.NodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "shard0.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		sf.NodeRepo.Create(&domain.Node{
			ID:        "node2",
			ShardID:   1,
			Hostname:  "shard1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		defer gock.Off()

		for i := 0; i < 5; i++ {
			url := domain.URL(fmt.Sprintf("http://example.com/about%d", i))
			shardID, err := sf.ShardService.GetURLShardID(url)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			job := &domain.CrawlJob{
				PageID:  util.PageID(url),
				URL:     url,
				ShardID: shardID,
				Status:  domain.CrawlStatusPending,
			}

			err = sf.CrawlJobRepo.Save(job)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}

		for i := 0; i < 5; i++ {
			url := domain.URL(fmt.Sprintf("http://example.net/about%d", i))
			shardID, err := sf.ShardService.GetURLShardID(url)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			job := &domain.CrawlJob{
				PageID:  util.PageID(url),
				URL:     url,
				ShardID: shardID,
				Status:  domain.CrawlStatusPending,
			}

			err = sf.CrawlJobRepo.Save(job)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}

		gock.New("http://shard0.cluster.com:8080").
			Post("/crawl").
			Times(1).
			Reply(200).
			JSON(dto.CrawlResponse{
				Links: []string{
					"http://example.com/1",
					"http://example.com/2",
				},
			})

		gock.New("http://shard1.cluster.com:8080").
			Post("/crawl").
			Times(1).
			Reply(200).
			JSON(dto.CrawlResponse{
				Links: []string{
					"http://example.net/1",
					"http://example.net/2",
				},
			})

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		err := crawlService.RunCrawlProcess(ctx)

		if err != context.DeadlineExceeded {
			t.Errorf("expected got context deadline exceeded %v", err)
		}

		var pendingCount int
		var completedCount int

		for i := 0; i < 5; i++ {
			url := domain.URL(fmt.Sprintf("http://example.com/about%d", i))
			job, err := sf.CrawlJobRepo.Get(util.PageID(url))

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if job.Status == domain.CrawlStatusPending {
				pendingCount++
			}

			if job.Status == domain.CrawlStatusCompleted {
				completedCount++
			}
		}

		if pendingCount != 4 {
			t.Errorf("expected 4 pending jobs, got %v", pendingCount)
		}

		if completedCount != 1 {
			t.Errorf("expected 1 completed job, got %v", completedCount)
		}

		pendingCount = 0
		completedCount = 0

		for i := 0; i < 5; i++ {
			url := domain.URL(fmt.Sprintf("http://example.net/about%d", i))
			job, err := sf.CrawlJobRepo.Get(util.PageID(url))

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if job.Status == domain.CrawlStatusPending {
				pendingCount++
			}

			if job.Status == domain.CrawlStatusCompleted {
				completedCount++
			}
		}

		if pendingCount != 4 {
			t.Errorf("expected 4 pending jobs, got %v", pendingCount)
		}

		if completedCount != 1 {
			t.Errorf("expected 1 completed job, got %v", completedCount)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})
}

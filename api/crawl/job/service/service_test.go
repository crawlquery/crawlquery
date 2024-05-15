package service_test

import (
	"crawlquery/api/crawl/job/repository/mem"
	"crawlquery/api/crawl/job/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"fmt"
	"time"

	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	"errors"
	"testing"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	t.Run("can create a job", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, testutil.NewTestLogger())
		url := "http://example.com"

		// Act
		job, err := svc.Create(url)

		// Assert
		if err != nil {
			t.Errorf("Error adding job: %v", err)
		}

		if job.ID == "" {
			t.Errorf("Expected ID to be set")
		}

		if job.URL != url {
			t.Errorf("Expected URL to be %s, got %s", url, job.URL)
		}

		if job.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set")
		}
	})

	t.Run("validates url", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, testutil.NewTestLogger())
		url := "x123!"

		// Act
		job, err := svc.Create(url)

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if job != nil {
			t.Errorf("Expected job to be nil")
		}
	})

	t.Run("handles repository error", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, testutil.NewTestLogger())
		expectErr := errors.New("db locked")
		repo.ForceError(expectErr)

		// Act
		job, err := svc.Create("http://example.com")

		// Assert
		if err != domain.ErrInternalError {
			t.Errorf("Expected error, got nil")
		}

		if job != nil {
			t.Errorf("Expected job to be nil")
		}
	})
}

func TestProcessCrawlJobs(t *testing.T) {
	t.Run("can process crawl jobs", func(t *testing.T) {

		nodeRepo := nodeRepo.NewRepository()
		nodeService := nodeService.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		nodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		nodeRepo.Create(&domain.Node{
			ID:        "node2",
			ShardID:   1,
			Hostname:  "node2.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		shardRepo := shardRepo.NewRepository()
		shardService := shardService.NewService(shardRepo, testutil.NewTestLogger())

		shardRepo.Create(&domain.Shard{
			ID: 0,
		})

		shardRepo.Create(&domain.Shard{
			ID: 1,
		})

		defer gock.Off()

		url := "http://example.com"

		pageID := util.PageID(url)

		responseJson := fmt.Sprintf(`{"page_id":"%s","url":"%s"}`, pageID, url)

		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			JSON(responseJson).
			Reply(200)

		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, shardService, nodeService, testutil.NewTestLogger())
		job, _ := svc.Create(url)

		// Act
		go svc.ProcessCrawlJobs()

		// Wait for the job to be processed
		time.Sleep(100 * time.Millisecond)

		// Assert
		job, err := repo.Get(job.ID)

		if err != nil {
			t.Errorf("Error getting job: %v", err)
		}

		if job == nil {
			t.Fatalf("Expected job to be set")
		}

		if job.LastCrawledAt.Time.IsZero() {
			t.Errorf("Expected LastCrawledAt to be set")
		}

		if job.FailedReason.String != "" {
			t.Errorf("Expected FailedReason to be empty")
		}

		if !job.BackoffUntil.Time.IsZero() {
			t.Errorf("Expected BackoffUntil to be empty")
		}
	})

	t.Run("handles repository error", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, testutil.NewTestLogger())
		expectErr := errors.New("db locked")
		repo.ForceError(expectErr)

		// Act
		go svc.ProcessCrawlJobs()
	})
}

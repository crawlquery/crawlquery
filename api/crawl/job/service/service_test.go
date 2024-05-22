package service_test

import (
	"crawlquery/api/crawl/job/repository/mem"
	"crawlquery/api/crawl/job/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"database/sql"
	"fmt"
	"strings"
	"time"

	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	resRepo "crawlquery/api/crawl/restriction/repository/mem"
	resService "crawlquery/api/crawl/restriction/service"

	"errors"
	"testing"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	t.Run("can create a job", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, testutil.NewTestLogger())
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
		svc := service.NewService(repo, nil, nil, nil, testutil.NewTestLogger())
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
		svc := service.NewService(repo, nil, nil, nil, testutil.NewTestLogger())
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

		resRepo := resRepo.NewRepository()
		resService := resService.NewService(resRepo)

		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, shardService, nodeService, resService, testutil.NewTestLogger())
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

		res, err := resRepo.Get("example.com")

		if err != nil {
			t.Errorf("Error getting restriction: %v", err)
		}

		if res == nil {
			t.Fatalf("Expected restriction to be set")
		}

		if !res.Until.Valid {
			t.Errorf("Expected restriction until to be set")
		}

		// Assert that the restriction is for 5 minutes
		if res.Until.Time.Round(time.Minute) != time.Now().Add(5*time.Minute).Round(time.Minute) {
			t.Errorf("Expected restriction until to be 5 minutes from now")
		}
	})

	t.Run("cannot process crawl job if restriction exists", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		resService := resService.NewService(resRepo)

		resRepo.Set(&domain.CrawlRestriction{
			Domain: "example.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		})

		url := "http://example.com/homepage"
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, resService, testutil.NewTestLogger())
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

		if !job.LastCrawledAt.Time.IsZero() {
			t.Errorf("Expected LastCrawledAt to be empty")
		}

		if !strings.Contains(job.FailedReason.String, "domain is restricted until") {
			t.Errorf("Expected FailedReason to be set, got %s", job.FailedReason.String)
		}
	})

	t.Run("handles repository error", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, testutil.NewTestLogger())
		expectErr := errors.New("db locked")
		repo.ForceError(expectErr)

		// Act
		go svc.ProcessCrawlJobs()
	})
}

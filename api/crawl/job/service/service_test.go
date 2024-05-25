package service_test

import (
	"crawlquery/api/crawl/job/repository/mem"
	"crawlquery/api/crawl/job/service"
	"crawlquery/api/domain"
	"crawlquery/node/dto"
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

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"

	indexJobRepo "crawlquery/api/index/job/repository/mem"
	indexJobService "crawlquery/api/index/job/service"

	"errors"
	"testing"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	t.Run("can create a job", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
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
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
		url := "notaurl"

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

	t.Run("normalizes url", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
		url := "http://example.com?utm_source=google&utm_medium=cpc&utm_campaign=summer-sale"

		// Act
		job, err := svc.Create(url)

		// Assert
		if err != nil {
			t.Errorf("Error adding job: %v", err)
		}

		if job.URL != "http://example.com" {
			t.Errorf("Expected URL to be %s, got %s", url, job.URL)
		}
	})

	t.Run("handles repository error", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
		expectErr := errors.New("db locked")
		repo.ForceError(expectErr)

		// Act
		job, err := svc.Create("http://example.com")

		// Assert
		if err == nil {
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

		crawlResponse := &dto.CrawlResponse{
			Page: &dto.Page{
				ID:   pageID,
				URL:  url,
				Hash: "hash123",
			},
		}
		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			JSON(responseJson).
			Reply(200).
			JSON(crawlResponse)

		resRepo := resRepo.NewRepository()
		resService := resService.NewService(resRepo, testutil.NewTestLogger())

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil, testutil.NewTestLogger())

		indexJobRepo := indexJobRepo.NewRepository()
		indexJobService := indexJobService.NewService(indexJobRepo, testutil.NewTestLogger())

		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, shardService, nodeService, resService, pageService, indexJobService, testutil.NewTestLogger())
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
		if time.Until(res.Until.Time).Round(time.Minute) != time.Minute*5 {
			t.Errorf("Expected restriction until to be 5 minutes from now")
		}

		indexJob, err := indexJobRepo.GetByPageID(pageID)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}

		if indexJob == nil {
			t.Fatalf("Expected index job to be set")
		}
	})

	t.Run("cannot process crawl job if restriction exists", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		resService := resService.NewService(resRepo, testutil.NewTestLogger())

		resRepo.Set(&domain.CrawlRestriction{
			Domain: "example.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		})

		url := "http://example.com/homepage"
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, resService, nil, nil, testutil.NewTestLogger())
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

	t.Run("creates page if crawl is successful", func(t *testing.T) {
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

		crawlResponse := &dto.CrawlResponse{
			&dto.Page{
				ID:   pageID,
				URL:  url,
				Hash: "hash123",
			},
		}
		gock.New("http://node1.cluster.com:8080").
			Post("/crawl").
			JSON(responseJson).
			Reply(200).
			JSON(crawlResponse)

		resRepo := resRepo.NewRepository()
		resService := resService.NewService(resRepo, testutil.NewTestLogger())

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil, testutil.NewTestLogger())

		indexJobRepo := indexJobRepo.NewRepository()
		indexJobService := indexJobService.NewService(indexJobRepo, testutil.NewTestLogger())

		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, shardService, nodeService, resService, pageService, indexJobService, testutil.NewTestLogger())
		svc.Create(url)

		// Act
		go svc.ProcessCrawlJobs()

		// Wait for the job to be processed
		time.Sleep(100 * time.Millisecond)

		pageCheck, err := pageRepo.Get(pageID)

		if err != nil {
			t.Errorf("Error getting page: %v", err)
		}

		if pageCheck == nil {
			t.Fatalf("Expected page to be set")
		}

		if pageCheck.ID != pageID {
			t.Errorf("Expected page ID to be %s, got %s", pageID, pageCheck.ID)
		}

		if pageCheck.ShardID != 0 {
			t.Errorf("Expected page ShardID to be 0, got %d", pageCheck.ShardID)
		}

		if pageCheck.Hash != crawlResponse.Page.Hash {
			t.Errorf("Expected page Hash to be %s, got %s", crawlResponse.Page.Hash, pageCheck.Hash)
		}
	})

	t.Run("handles repository error", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
		expectErr := errors.New("db locked")
		repo.ForceError(expectErr)

		// Act
		go svc.ProcessCrawlJobs()
	})
}

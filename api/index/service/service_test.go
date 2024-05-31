package service_test

import (
	"context"
	"crawlquery/api/domain"
	indexJobRepo "crawlquery/api/index/job/repository/mem"
	indexLogRepo "crawlquery/api/index/log/repository/mem"
	indexService "crawlquery/api/index/service"
	"crawlquery/api/testfactory"
	"crawlquery/node/dto"
	"fmt"
	"math/rand"
	"time"

	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"

	"github.com/h2non/gock"
)

func TestCreateJob(t *testing.T) {
	t.Run("can create index job", func(t *testing.T) {
		indexJobRepo := indexJobRepo.NewRepository()
		indexLogRepo := indexLogRepo.NewRepository()
		indexService := indexService.NewService(
			indexService.WithIndexLogRepo(indexLogRepo),
			indexService.WithIndexJobRepo(indexJobRepo),
			indexService.WithLogger(testutil.NewTestLogger()),
		)

		err := indexService.CreateJob("page1", 0)

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		jobs, err := indexJobRepo.ListByStatus(10, domain.IndexStatusPending)

		if err != nil {
			t.Errorf("Error listing index jobs: %v", err)
		}

		if len(jobs) != 1 {
			t.Errorf("Expected 1 job, got %v", len(jobs))
		}

		if jobs[0].PageID != "page1" {
			t.Errorf("Expected job ID to be page1, got %s", jobs[0].PageID)
		}

		if jobs[0].ShardID != 0 {
			t.Errorf("Expected shard ID to be 0, got %d", jobs[0].ShardID)
		}
	})

	t.Run("returns error if job already exists", func(t *testing.T) {
		indexJobRepo := indexJobRepo.NewRepository()
		indexService := indexService.NewService(
			indexService.WithIndexJobRepo(indexJobRepo),
			indexService.WithLogger(testutil.NewTestLogger()),
		)

		job := &domain.IndexJob{
			PageID: "job1",
		}

		indexJobRepo.Save(job)

		err := indexService.CreateJob("job1", 0)

		if err != domain.ErrIndexJobAlreadyExists {
			t.Errorf("Expected ErrIndexJobAlreadyExists, got %v", err)
		}
	})
}

func TestRunIndexProcess(t *testing.T) {

	defer gock.Off()
	t.Run("should process index jobs with 10 workers and 100 index jobs", func(t *testing.T) {
		sf := testfactory.NewServiceFactory()

		indexService := indexService.NewService(
			indexService.WithEventService(sf.EventService),
			indexService.WithIndexJobRepo(sf.IndexJobRepo),
			indexService.WithNodeService(sf.NodeService),
			indexService.WithIndexLogRepo(sf.IndexLogRepo),
			indexService.WithLogger(testutil.NewTestLogger()),
			indexService.WithWorkers(10),
			indexService.WithMaxQueueSize(100),
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

		shardJobs := make(map[domain.ShardID][]*domain.IndexJob)

		for i := 0; i < 100; i++ {
			url := domain.URL(fmt.Sprintf("http://example%d.com/about", i))
			shardID, err := sf.ShardService.GetURLShardID(url)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			job := &domain.IndexJob{
				PageID:  util.PageID(url),
				ShardID: shardID,
				Status:  domain.IndexStatusPending,
			}

			err = sf.IndexJobRepo.Save(job)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			shardJobs[shardID] = append(shardJobs[shardID], job)

			hostname := fmt.Sprintf("shard%d.cluster.com", shardID)

			if i == 50 {
				gock.New(hostname).
					Post(fmt.Sprintf("/pages/%s/index", job.PageID)).
					Reply(200).
					JSON(dto.ErrorResponse{
						Error: "invalid html",
					})
			} else {
				gock.New(hostname).
					Post(fmt.Sprintf("/pages/%s/index", job.PageID)).
					Reply(200).
					JSON(dto.IndexResponse{
						Success: true,
						Message: "indexed",
					})
			}

		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		err := indexService.RunIndexProcess(ctx)

		if err != context.DeadlineExceeded {
			t.Errorf("expected got context deadline exceeded %v", err)
		}

		var failedCount int
		var completedCount int

		for _, jobs := range shardJobs {
			for _, job := range jobs {
				updatedJob, err := sf.IndexJobRepo.Get(job.PageID)

				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if updatedJob.Status == domain.IndexStatusFailed {
					failedCount++
				}

				if updatedJob.Status == domain.IndexStatusCompleted {
					completedCount++
				}
			}
		}

		if failedCount != 1 {
			t.Errorf("expected 1 failed job, got %v", failedCount)
		}

		if completedCount != 99 {
			t.Errorf("expected 99 completed jobs, got %v", completedCount)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("should try 3 times to process job", func(t *testing.T) {
		rand.Seed(0)
		sf := testfactory.NewServiceFactory()

		indexService := indexService.NewService(
			indexService.WithEventService(sf.EventService),
			indexService.WithIndexJobRepo(sf.IndexJobRepo),
			indexService.WithNodeService(sf.NodeService),
			indexService.WithIndexLogRepo(sf.IndexLogRepo),
			indexService.WithLogger(testutil.NewTestLogger()),
			indexService.WithWorkers(10),
			indexService.WithMaxQueueSize(100),
		)

		sf.ShardRepo.Create(&domain.Shard{ID: 0})
		sf.ShardRepo.Create(&domain.Shard{ID: 1})

		sf.NodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		sf.NodeRepo.Create(&domain.Node{
			ID:        "node2",
			ShardID:   0,
			Hostname:  "node2.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		sf.NodeRepo.Create(&domain.Node{
			ID:        "node3",
			ShardID:   0,
			Hostname:  "node3.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		url := domain.URL("http://example.com/about")
		shardID, err := sf.ShardService.GetURLShardID(url)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		job := &domain.IndexJob{
			PageID:  util.PageID(url),
			ShardID: shardID,
			Status:  domain.IndexStatusPending,
		}

		err = sf.IndexJobRepo.Save(job)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		gock.New("node1.cluster.com").
			Post(fmt.Sprintf("/pages/%s/index", job.PageID)).
			Reply(500).
			JSON(dto.ErrorResponse{
				Error: "internal server error",
			})

		gock.New("node2.cluster.com").
			Post(fmt.Sprintf("/pages/%s/index", job.PageID)).
			Reply(500).
			JSON(dto.ErrorResponse{
				Error: "internal server error",
			})

		gock.New("node3.cluster.com").
			Post(fmt.Sprintf("/pages/%s/index", job.PageID)).
			Reply(200).
			JSON(dto.IndexResponse{
				Success: true,
				Message: "indexed",
			})

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		err = indexService.RunIndexProcess(ctx)

		if err != context.DeadlineExceeded {
			t.Errorf("expected got context deadline exceeded %v", err)
		}

		updatedJob, err := sf.IndexJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if updatedJob.Status != domain.IndexStatusCompleted {
			t.Errorf("expected job to be completed, got %v", updatedJob.Status)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})
}

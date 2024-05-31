package service_test

import (
	"crawlquery/api/domain"
	indexJobRepo "crawlquery/api/index/job/repository/mem"
	indexLogRepo "crawlquery/api/index/log/repository/mem"
	indexService "crawlquery/api/index/service"

	"crawlquery/pkg/testutil"
	"testing"
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

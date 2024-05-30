package service_test

import (
	"crawlquery/api/domain"
	indexJobRepo "crawlquery/api/index/job/repository/mem"
	indexService "crawlquery/api/index/service"

	"crawlquery/pkg/testutil"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create index job", func(t *testing.T) {
		indexJobRepo := indexJobRepo.NewRepository()
		indexService := indexService.NewService(
			indexService.WithIndexJobRepo(indexJobRepo),
			indexService.WithLogger(testutil.NewTestLogger()),
		)

		job, err := indexService.Create("job1")

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		if job.PageID != "job1" {
			t.Errorf("Expected job ID to be job1, got %s", job.PageID)
		}

		if job.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set")
		}

		checkJob, err := indexJobRepo.Get(job.PageID)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}

		if checkJob.PageID != job.PageID {
			t.Errorf("Expected job ID to be %s, got %s", job.PageID, checkJob.PageID)
		}

		if checkJob.CreatedAt != job.CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, got %v", job.CreatedAt, checkJob.CreatedAt)
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

		_, err := indexService.Create("job1")

		if err != domain.ErrIndexJobAlreadyExists {
			t.Errorf("Expected ErrIndexJobAlreadyExists, got %v", err)
		}
	})
}

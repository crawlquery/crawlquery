package service_test

import (
	"crawlquery/api/domain"
	indexJobRepo "crawlquery/api/index/job/repository/mem"
	indexJobService "crawlquery/api/index/job/service"

	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	t.Run("can create index job", func(t *testing.T) {
		indexJobRepo := indexJobRepo.NewRepository()
		indexJobService := indexJobService.NewService(indexJobRepo, nil, nil, testutil.NewTestLogger())

		job, err := indexJobService.Create("job1")

		if err != nil {
			t.Errorf("Error creating index job: %v", err)
		}

		if job.ID == "" {
			t.Errorf("Expected job ID to be set")
		}

		if uuid.Validate(job.ID) != nil {
			t.Errorf("Expected job ID to be a valid UUID")
		}

		if job.PageID != "job1" {
			t.Errorf("Expected job ID to be job1, got %s", job.PageID)
		}

		if job.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set")
		}

		checkJob, err := indexJobRepo.Get(job.ID)

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
		indexJobService := indexJobService.NewService(indexJobRepo, nil, nil, testutil.NewTestLogger())

		job := &domain.IndexJob{
			PageID: "job1",
		}

		indexJobRepo.Create(job)

		_, err := indexJobService.Create("job1")

		if err != domain.ErrIndexJobAlreadyExists {
			t.Errorf("Expected ErrIndexJobAlreadyExists, got %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("can get index job", func(t *testing.T) {
		indexJobRepo := indexJobRepo.NewRepository()
		indexJobService := indexJobService.NewService(indexJobRepo, nil, nil, testutil.NewTestLogger())

		job := &domain.IndexJob{
			PageID: "job1",
		}

		indexJobRepo.Create(job)

		result, err := indexJobService.Get(job.ID)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}

		if result.PageID != job.PageID {
			t.Errorf("Expected job ID to be %s, got %s", job.PageID, result.PageID)
		}
	})
}

func TestNext(t *testing.T) {
	t.Run("can get next index job", func(t *testing.T) {
		indexJobRepo := indexJobRepo.NewRepository()
		indexJobService := indexJobService.NewService(indexJobRepo, nil, nil, testutil.NewTestLogger())

		job := &domain.IndexJob{
			PageID: "job1",
		}

		indexJobRepo.Create(job)

		result, err := indexJobService.Next()

		if err != nil {
			t.Errorf("Error getting next index job: %v", err)
		}

		if result.PageID != job.PageID {
			t.Errorf("Expected job ID to be %s, got %s", job.PageID, result.PageID)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("can update index job", func(t *testing.T) {
		indexJobRepo := indexJobRepo.NewRepository()
		indexJobService := indexJobService.NewService(indexJobRepo, nil, nil, testutil.NewTestLogger())

		job := &domain.IndexJob{
			PageID: "job1",
		}

		indexJobRepo.Create(job)

		job.PageID = "job2"

		err := indexJobService.Update(job)

		if err != nil {
			t.Errorf("Error updating index job: %v", err)
		}

		result, err := indexJobRepo.Get(job.ID)

		if err != nil {
			t.Errorf("Error getting index job: %v", err)
		}

		if result.PageID != job.PageID {
			t.Errorf("Expected job ID to be %s, got %s", job.PageID, result.PageID)
		}
	})
}

func TestProcessIndexJobs(t *testing.T) {
	t.Run("can process index jobs", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/pages/page1/index").
			Reply(200).
			JSON(`{"success": true,"message": "Page indexed"}`)

		nodeRepo := nodeRepo.NewRepository()
		nodeService := nodeService.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		nodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil, testutil.NewTestLogger())

		pageRepo.Create(&domain.Page{
			ID:        "page1",
			ShardID:   0,
			CreatedAt: time.Now(),
		})

		indexJobRepo := indexJobRepo.NewRepository()
		indexJobService := indexJobService.NewService(indexJobRepo, pageService, nodeService, testutil.NewTestLogger())

		job := &domain.IndexJob{
			PageID: "page1",
		}

		indexJobRepo.Create(job)

		go indexJobService.ProcessIndexJobs()

		// Wait for the job to be processed
		time.Sleep(100 * time.Millisecond)

		// Assert
		job, err := indexJobRepo.Get(job.ID)

		if err != nil {
			t.Errorf("Error getting job: %v", err)
		}

		if job == nil {
			t.Fatalf("Expected job to be set")
		}

		if job.LastIndexedAt.Time.IsZero() {
			t.Errorf("Expected LastIndexedAt to be set")
		}
	})
}

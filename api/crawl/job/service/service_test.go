package service_test

import (
	"crawlquery/api/crawl/job/repository/mem"
	"crawlquery/api/crawl/job/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/testutil"
	"errors"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create a job", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())
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
		svc := service.NewService(repo, testutil.NewTestLogger())
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
		svc := service.NewService(repo, testutil.NewTestLogger())
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

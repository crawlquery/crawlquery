package service_test

import (
	"crawlquery/api/crawl/job/repository/mem"
	"crawlquery/api/crawl/service"
	"crawlquery/pkg/testutil"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestAddJob(t *testing.T) {
	t.Run("can add a job", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())
		url := "http://example.com"

		// Act
		err := svc.AddJob(url)

		// Assert
		if err != nil {
			t.Errorf("Error adding job: %v", err)
		}

		first, err := repo.First()
		if err != nil {
			t.Errorf("Error getting first job: %v", err)
		}

		if first.URL != url {
			t.Errorf("Expected URL to be %s, got %s", url, first.URL)
		}

		if _, err := uuid.Parse(first.ID); err != nil {
			t.Errorf("Expected ID to be a valid UUID, got %s", first.ID)
		}
	})

	t.Run("validates url", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())
		url := "x123!"

		// Act
		err := svc.AddJob(url)

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("handles repository error", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())
		expectErr := errors.New("db locked")
		repo.ForceError(expectErr)

		// Act
		err := svc.AddJob("http://example.com")

		// Assert
		if err != expectErr {
			t.Errorf("Expected error, got nil")
		}
	})
}

package domain_test

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"strings"
	"testing"
	"time"
)

func TestCrawlJobValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cj := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "http://example.com",
			CreatedAt: time.Now(),
		}

		err := cj.Validate()

		if err != nil {
			t.Errorf("Expected crawl job to be valid, got error: %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cj := &domain.CrawlJob{
			ID:        "abs",
			URL:       "http://example.com",
			CreatedAt: time.Now(),
		}

		err := cj.Validate()

		if err == nil {
			t.Errorf("Expected crawl job to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "CrawlJob.ID") {
			t.Errorf("Expected error to contain 'CrawlJob.ID', got %v", err)
		}
	})

	t.Run("invalid url", func(t *testing.T) {
		cj := &domain.CrawlJob{
			ID:        util.UUID(),
			URL:       "example.com",
			CreatedAt: time.Now(),
		}

		err := cj.Validate()

		if err == nil {
			t.Errorf("Expected crawl job to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "CrawlJob.URL") {
			t.Errorf("Expected error to contain 'CrawlJob.URL', got %v", err)
		}
	})
}

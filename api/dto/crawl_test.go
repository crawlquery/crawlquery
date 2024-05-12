package dto_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"testing"
	"time"
)

func TestNewCreateCrawlJobResponse(t *testing.T) {
	t.Run("should return correct CreateCrawlResponse from Crawl", func(t *testing.T) {
		// given
		c := &domain.CrawlJob{
			ID:        "1",
			URL:       "http://example.com",
			CreatedAt: time.Now(),
		}

		r := dto.NewCreateCrawlJobResponse(c)

		// then
		if r.CrawlJob.ID != c.ID {
			t.Errorf("Expected ID to be %s, got %s", c.ID, r.CrawlJob.ID)
		}

		if r.CrawlJob.URL != c.URL {
			t.Errorf("Expected URL to be %s, got %s", c.URL, r.CrawlJob.URL)
		}

		if r.CrawlJob.CreatedAt != c.CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, got %v", c.CreatedAt, r.CrawlJob.CreatedAt)
		}

	})
}

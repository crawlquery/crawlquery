package dto_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"testing"
	"time"
)

func TestCreateCrawlJobRequestToJSON(t *testing.T) {
	t.Run("should return correct JSON", func(t *testing.T) {
		// given
		r := &dto.CreateCrawlJobRequest{
			URL: "http://example.com",
		}

		// when
		b, err := r.ToJSON()

		// then
		if err != nil {
			t.Errorf("Expected error to be nil, got %v", err)
		}

		expected := `{"url":"http://example.com"}`
		if string(b) != expected {
			t.Errorf("Expected JSON to be %s, got %s", expected, string(b))
		}
	})
}

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

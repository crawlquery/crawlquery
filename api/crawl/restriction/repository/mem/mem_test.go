package mem

import (
	"crawlquery/api/domain"
	"database/sql"
	"testing"
	"time"
)

func TestRepositoryGet(t *testing.T) {
	t.Run("returns crawl restriction if it exists", func(t *testing.T) {
		repo := NewRepository()

		restriction := &domain.CrawlRestriction{
			Domain: "example.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		}

		repo.Set(restriction)

		got, err := repo.Get("example.com")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if got != restriction {
			t.Errorf("Expected %v, got %v", restriction, got)
		}
	})

	t.Run("returns error if crawl restriction does not exist", func(t *testing.T) {
		repo := NewRepository()

		_, err := repo.Get("example.com")

		if err != domain.ErrCrawlRestrictionNotFound {
			t.Errorf("Expected %v, got %v", domain.ErrCrawlRestrictionNotFound, err)
		}
	})
}

func TestRepositorySet(t *testing.T) {
	t.Run("sets crawl restriction", func(t *testing.T) {
		repo := NewRepository()

		restriction := &domain.CrawlRestriction{
			Domain: "example.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		}

		err := repo.Set(restriction)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		got, err := repo.Get("example.com")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if got != restriction {
			t.Errorf("Expected %v, got %v", restriction, got)
		}
	})
}

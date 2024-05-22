package service_test

import (
	resRepo "crawlquery/api/crawl/restriction/repository/mem"
	"crawlquery/api/crawl/restriction/service"
	"crawlquery/api/domain"
	"database/sql"
	"testing"
	"time"
)

func TestHasRestriction(t *testing.T) {
	t.Run("returns true if there is a restriction", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		service := service.NewService(resRepo)

		resRepo.Set(&domain.CrawlRestriction{
			Domain: "http://example.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		})

		if !service.HasRestriction("http://example.com") {
			t.Errorf("expected HasRestriction to return true, got false")
		}
	})

	t.Run("returns false if there is no restriction", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		service := service.NewService(resRepo)

		if service.HasRestriction("http://example.com") {
			t.Errorf("expected HasRestriction to return false, got true")
		}
	})
}

func TestRestrict(t *testing.T) {
	t.Run("returns error if restriction already exists", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		service := service.NewService(resRepo)

		resRepo.Set(&domain.CrawlRestriction{
			Domain: "http://example.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		})

		err := service.Restrict("http://example.com")

		if err != domain.ErrCrawlRestrictionAlreadyExists {
			t.Errorf("expected ErrCrawlRestrictionAlreadyExists, got %v", err)
		}
	})

	t.Run("sets a restriction", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		service := service.NewService(resRepo)

		err := service.Restrict("http://example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !service.HasRestriction("http://example.com") {
			t.Errorf("expected HasRestriction to return true, got false")
		}
	})
}

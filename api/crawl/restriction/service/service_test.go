package service_test

import (
	resRepo "crawlquery/api/crawl/restriction/repository/mem"
	"crawlquery/api/crawl/restriction/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/testutil"
	"database/sql"
	"testing"
	"time"
)

func TestGetRestriction(t *testing.T) {
	t.Run("returns true if there is a restriction", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		service := service.NewService(resRepo, testutil.NewTestLogger())
		res := &domain.CrawlRestriction{
			Domain: "http://example.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		}
		resRepo.Set(res)
		restricted, until := service.GetRestriction("http://example.com")
		if !restricted {
			t.Errorf("expected GetRestriction to return true, got false")
		}

		if until == nil {
			t.Errorf("expected until to be non-nil, got nil")
		}

		if *until != res.Until.Time {
			t.Errorf("expected until to be %v, got %v", res.Until.Time, until)
		}
	})

	t.Run("returns false if there is no restriction", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		service := service.NewService(resRepo, testutil.NewTestLogger())

		if restricted, _ := service.GetRestriction("http://example.com"); restricted {
			t.Errorf("expected GetRestriction to return false, got true")
		}
	})
}

func TestRestrict(t *testing.T) {
	t.Run("returns error if restriction already exists", func(t *testing.T) {
		resRepo := resRepo.NewRepository()
		service := service.NewService(resRepo, testutil.NewTestLogger())

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
		service := service.NewService(resRepo, testutil.NewTestLogger())

		err := service.Restrict("http://example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		res, err := resRepo.Get("http://example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if res.Domain != "http://example.com" {
			t.Errorf("expected domain to be 'http://example.com', got %v", res.Domain)
		}

		if time.Until(res.Until.Time).Round(time.Second) != time.Second*20 {
			t.Errorf("expected until to be 20 seconds from now, got %v", time.Until(res.Until.Time))
		}
	})
}

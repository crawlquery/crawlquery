package service_test

import (
	"crawlquery/pkg/testutil"
	"testing"

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"
)

func TestGet(t *testing.T) {
	t.Run("returns page", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, testutil.NewTestLogger())

		page, err := pageService.Create("pageID", 0)

		if err != nil {
			t.Fatalf("error creating page: %v", err)
		}

		got, err := pageService.Get(page.ID)

		if err != nil {
			t.Fatalf("error getting page: %v", err)
		}

		if got.ID != page.ID {
			t.Errorf("got page ID %s, want %s", got.ID, page.ID)
		}
	})

	t.Run("returns err if page not found", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, testutil.NewTestLogger())

		_, err := pageService.Get("pageID")

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("creates page", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, testutil.NewTestLogger())

		page, err := pageService.Create("pageID", 0)

		if err != nil {
			t.Fatalf("error creating page: %v", err)
		}

		if page.ID != "pageID" {
			t.Errorf("got page ID %s, want pageID", page.ID)
		}

		if page.ShardID != 0 {
			t.Errorf("got page ShardID %d, want 0", page.ShardID)
		}
	})

	t.Run("returns error if page already exists", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, testutil.NewTestLogger())

		_, err := pageService.Create("pageID", 0)

		if err != nil {
			t.Fatalf("error creating page: %v", err)
		}

		_, err = pageService.Create("pageID", 0)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

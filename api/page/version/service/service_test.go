package service_test

import (
	"crawlquery/api/domain"
	pageVersionRepo "crawlquery/api/page/version/repository/mem"
	pageVersionService "crawlquery/api/page/version/service"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Run("returns a page version", func(t *testing.T) {
		pageVersionRepo := pageVersionRepo.NewRepository()
		pageVersionService := pageVersionService.NewService(
			pageVersionService.WithVersionRepo(pageVersionRepo),
		)

		pageVersion := &domain.PageVersion{
			ID:          "id",
			PageID:      "page1",
			ContentHash: "hash",
			CreatedAt:   time.Now(),
		}

		pageVersionRepo.Create(pageVersion)

		got, err := pageVersionService.Get("id")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if got.ID != pageVersion.ID {
			t.Fatalf("Expected %s, got %s", pageVersion.ID, got.ID)
		}

		if got.PageID != pageVersion.PageID {
			t.Fatalf("Expected %s, got %s", pageVersion.PageID, got.PageID)
		}

		if got.ContentHash != pageVersion.ContentHash {
			t.Fatalf("Expected %s, got %s", pageVersion.ContentHash, got.ContentHash)
		}

		if !got.CreatedAt.Equal(pageVersion.CreatedAt) {
			t.Fatalf("Expected %v, got %v", pageVersion.CreatedAt, got.CreatedAt)
		}

	})
}

func TestListByPageID(t *testing.T) {
	t.Run("returns page versions", func(t *testing.T) {
		pageVersionRepo := pageVersionRepo.NewRepository()
		pageVersionService := pageVersionService.NewService(
			pageVersionService.WithVersionRepo(pageVersionRepo),
		)

		pageID := domain.PageID("page1")

		pageVersion1 := &domain.PageVersion{
			ID:          "id1",
			PageID:      pageID,
			ContentHash: "hash1",
			CreatedAt:   time.Now(),
		}

		pageVersion2 := &domain.PageVersion{
			ID:          "id2",
			PageID:      pageID,
			ContentHash: "hash2",
			CreatedAt:   time.Now(),
		}

		pageVersion3 := &domain.PageVersion{
			ID:          "id3",
			PageID:      "page2",
			ContentHash: "hash3",
			CreatedAt:   time.Now(),
		}

		pageVersionRepo.Create(pageVersion1)
		pageVersionRepo.Create(pageVersion2)
		pageVersionRepo.Create(pageVersion3)

		got, err := pageVersionService.ListByPageID(pageID)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("Expected 2 page versions, got %d", len(got))
		}

		if got[0].ID != pageVersion1.ID && got[1].ID != pageVersion1.ID {
			t.Fatalf("Expected %s, got %s", pageVersion1.ID, got[0].ID)
		}

	})
}

func TestCreate(t *testing.T) {
	t.Run("creates a page version", func(t *testing.T) {
		pageVersionRepo := pageVersionRepo.NewRepository()
		pageVersionService := pageVersionService.NewService(
			pageVersionService.WithVersionRepo(pageVersionRepo),
		)

		now := time.Now()

		pageVersion := &domain.PageVersion{
			PageID:      "page1",
			ContentHash: "hash",
		}

		got, err := pageVersionService.Create(pageVersion.PageID, pageVersion.ContentHash)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if got.ID == "" {
			t.Fatalf("Expected page version to have an ID")
		}

		if got.PageID != pageVersion.PageID {
			t.Fatalf("Expected %s, got %s", pageVersion.PageID, got.PageID)
		}

		if got.ContentHash != pageVersion.ContentHash {
			t.Fatalf("Expected %s, got %s", pageVersion.ContentHash, got.ContentHash)
		}

		if !got.CreatedAt.Round(time.Second).Equal(now.Round(time.Second)) {
			t.Fatalf("Expected %v, got %v", now, got.CreatedAt)
		}
	})
}

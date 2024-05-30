package mem

import (
	"crawlquery/api/domain"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("returns a page version", func(t *testing.T) {
		repo := NewRepository()
		pageVersion := &domain.PageVersion{
			ID:          "id",
			PageID:      "page1",
			ContentHash: "hash",
		}

		repo.Create(pageVersion)

		got, err := repo.Get("id")

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
	})
}

func TestListByPageID(t *testing.T) {
	t.Run("returns page versions", func(t *testing.T) {
		repo := NewRepository()
		pageID := domain.PageID("page1")

		pageVersion1 := &domain.PageVersion{
			ID:          "id1",
			PageID:      pageID,
			ContentHash: "hash1",
		}

		pageVersion2 := &domain.PageVersion{
			ID:          "id2",
			PageID:      pageID,
			ContentHash: "hash2",
		}

		repo.Create(pageVersion1)
		repo.Create(pageVersion2)

		got, err := repo.ListByPageID(pageID)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("Expected 2 page versions, got %d", len(got))
		}

		if got[0].ID != pageVersion1.ID {
			t.Fatalf("Expected %s, got %s", pageVersion1.ID, got[0].ID)
		}

		if got[1].ID != pageVersion2.ID {
			t.Fatalf("Expected %s, got %s", pageVersion2.ID, got[1].ID)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("creates a page version", func(t *testing.T) {
		repo := NewRepository()
		pageVersion := &domain.PageVersion{
			ID:          "id",
			PageID:      "page1",
			ContentHash: "hash",
		}

		err := repo.Create(pageVersion)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		got, err := repo.Get("id")

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
	})
}

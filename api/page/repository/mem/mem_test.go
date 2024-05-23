package mem

import (
	"crawlquery/api/domain"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("should save a page", func(t *testing.T) {
		repo := NewRepository()

		page := &domain.Page{
			ID:      "123",
			ShardID: 1,
		}

		err := repo.Create(page)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if repo.pages["123"] == nil {
			t.Errorf("expected page to be saved")
		}

		if repo.pages["123"].ID != "123" {
			t.Errorf("expected page ID to be 123, got %s", repo.pages["123"].ID)
		}

		if repo.pages["123"].ShardID != 1 {
			t.Errorf("expected page ShardID to be 1, got %d", repo.pages["123"].ShardID)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("should return a page", func(t *testing.T) {
		repo := NewRepository()

		page := &domain.Page{
			ID:      "123",
			ShardID: 1,
		}

		repo.pages["123"] = page

		res, err := repo.Get("123")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if res.ID != "123" {
			t.Errorf("expected page ID to be 123, got %s", res.ID)
		}

		if res.ShardID != 1 {
			t.Errorf("expected page ShardID to be 1, got %d", res.ShardID)
		}
	})

	t.Run("should return an error if page is not found", func(t *testing.T) {
		repo := NewRepository()

		_, err := repo.Get("123")

		if err != domain.ErrPageNotFound {
			t.Errorf("expected error: %v, got: %v", domain.ErrPageNotFound, err)
		}
	})
}

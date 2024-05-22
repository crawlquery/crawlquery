package service_test

import (
	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	"crawlquery/node/page/service"
	"testing"
)

func TestCreate(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	page, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	if page.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", page.ID)
	}

	if page.URL != "http://example.com" {
		t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", page.URL)
	}

	check, err := pageRepo.Get("1")

	if err != nil {
		t.Fatalf("Error getting page: %v", err)
	}

	if check.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", check.ID)
	}

	if check.URL != "http://example.com" {
		t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", check.URL)
	}
}

func TestGet(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	_, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	check, err := service.Get("1")

	if err != nil {
		t.Fatalf("Error getting page: %v", err)
	}

	if check.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", check.ID)
	}

	if check.URL != "http://example.com" {
		t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", check.URL)
	}
}

func TestGetAll(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	_, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	_, err = service.Create("2", "http://example2.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	pages, err := service.GetAll()

	if err != nil {
		t.Fatalf("Error getting pages: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("Expected 2 pages, got %d", len(pages))
	}
}

func TestUpdate(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	page, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	page.URL = "http://example2.com"

	err = service.Update(page)

	if err != nil {
		t.Fatalf("Error updating page: %v", err)
	}

	check, err := service.Get("1")

	if err != nil {
		t.Fatalf("Error getting page: %v", err)
	}

	if check.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", check.ID)
	}

	if check.URL != "http://example2.com" {
		t.Fatalf("Expected page URL to be 'http://example2.com', got '%s'", check.URL)
	}
}

func TestDelete(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	_, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	err = service.Delete("1")

	if err != nil {
		t.Fatalf("Error deleting page: %v", err)
	}

	_, err = service.Get("1")

	if err != domain.ErrPageNotFound {
		t.Fatalf("Expected ErrPageNotFound, got %v", err)
	}
}

func TestHash(t *testing.T) {

	t.Run("can create hash", func(t *testing.T) {
		repo := pageRepo.NewRepository()
		s := service.NewService(repo)
		_, err := s.Create("page1", "http://example.com")

		if err != nil {
			t.Fatalf("Error saving postings: %v", err)
		}

		hash, err := repo.GetHash("page1")

		if err != nil {
			t.Fatalf("Error getting hash: %v", err)
		}

		if hash == "" {
			t.Fatalf("Expected hash to not be empty")
		}
	})

	t.Run("can delete hash", func(t *testing.T) {
		repo := pageRepo.NewRepository()
		s := service.NewService(repo)
		_, err := s.Create("page1", "http://example.com")

		if err != nil {
			t.Fatalf("Error saving postings: %v", err)
		}

		hash, err := repo.GetHash("page1")

		if err != nil {
			t.Fatalf("Error getting hash: %v", err)
		}

		if hash == "" {
			t.Fatalf("Expected hash to not be empty")
		}
	})

	t.Run("can get hash of all pages", func(t *testing.T) {
		repo := pageRepo.NewRepository()
		s := service.NewService(repo)
		s.Create("page1", "http://example.com")
		s.Create("page2", "http://example2.com")

		hash1, err := s.Hash()

		if err != nil {
			t.Fatalf("Error getting hash: %v", err)
		}

		if hash1 == "" {
			t.Fatalf("Expected hash to not be empty")
		}

		_, err = s.Create("page3", "http://example3.com")

		if err != nil {
			t.Fatalf("Error saving postings: %v", err)
		}

		hash2, err := s.Hash()

		if err != nil {
			t.Fatalf("Error getting hash: %v", err)
		}

		if hash2 == "" {
			t.Fatalf("Expected hash to not be empty")
		}

		if hash1 == hash2 {
			t.Fatalf("Expected hashes to be different, got %s", hash1)
		}
	})
}

func TestJSON(t *testing.T) {
	repo := pageRepo.NewRepository()
	s := service.NewService(repo)
	_, err := s.Create("page1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	_, err = s.Create("page2", "http://example2.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	json, err := s.JSON()

	if err != nil {
		t.Fatalf("Error getting json: %v", err)
	}

	expected := `{"page1":{"id":"page1","url":"http://example.com","title":"","meta_description":""},"page2":{"id":"page2","url":"http://example2.com","title":"","meta_description":""}}`

	if string(json) != expected {
		t.Fatalf("Expected json to be %s, got %s", expected, json)
	}

}

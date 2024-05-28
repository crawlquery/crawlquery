package service_test

import (
	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	"crawlquery/node/page/service"
	peerService "crawlquery/node/peer/service"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"testing"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	t.Run("can create page", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		service := service.NewService(pageRepo, nil)

		page, err := service.Create("1", "http://example.com", "hash1")

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

		if check.Hash != "hash1" {
			t.Fatalf("Expected page Hash to be 'hash1', got '%s'", check.Hash)
		}
	})

	t.Run("broadasts event to peers", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		peerService.AddPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "peer1",
			Port:     8080,
			ShardID:  1,
		})

		expectedPage := &domain.Page{
			ID:   "1",
			URL:  "http://example.com",
			Hash: "hash1",
		}

		defer gock.Off()

		gock.New("http://peer1:8080").
			Post("/event").
			JSON(expectedPage).
			Reply(200)

		service := service.NewService(pageRepo, peerService)

		_, err := service.Create("1", "http://example.com", "hash1")

		if err != nil {
			t.Fatalf("Error saving page: %v", err)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

func TestCount(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo, nil)

	_, err := service.Create("1", "http://example.com", "hash1")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	count, err := service.Count()

	if err != nil {
		t.Fatalf("Error counting pages: %v", err)
	}

	if count != 1 {
		t.Fatalf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo, nil)

	_, err := service.Create("1", "http://example.com", "hash1")

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

	if check.Hash != "hash1" {
		t.Fatalf("Expected page Hash to be 'hash1', got '%s'", check.Hash)
	}
}

func TestGetAll(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo, nil)

	_, err := service.Create("1", "http://example.com", "hash1")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	_, err = service.Create("2", "http://example2.com", "hash2")

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
	t.Run("can update page", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		service := service.NewService(pageRepo, nil)

		page, err := service.Create("1", "http://example.com", "hash1")

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
	})

	t.Run("broadasts event to peers", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		peerService.AddPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "peer1",
			Port:     8080,
			ShardID:  1,
		})

		expectedPage := &domain.Page{
			ID:   "1",
			URL:  "http://example.com",
			Hash: "hash2",
		}

		defer gock.Off()

		gock.New("http://peer1:8080").
			Post("/event").
			JSON(expectedPage).
			Reply(200)

		service := service.NewService(pageRepo, peerService)

		err := pageRepo.Save(expectedPage.ID, expectedPage)

		if err != nil {
			t.Fatalf("Error saving page: %v", err)
		}

		page, err := service.Get("1")

		if err != nil {
			t.Fatalf("Error getting page: %v", err)
		}

		page.Hash = "hash2"

		err = service.Update(page)

		if err != nil {
			t.Fatalf("Error updating page: %v", err)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

func TestDelete(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo, nil)

	_, err := service.Create("1", "http://example.com", "hash1")

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
		s := service.NewService(repo, nil)
		_, err := s.Create("page1", "http://example.com", "hash1")

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
		s := service.NewService(repo, nil)
		_, err := s.Create("page1", "http://example.com", "hash1")

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
		s := service.NewService(repo, nil)
		s.Create("page1", "http://example.com", "hash1")
		s.Create("page2", "http://example2.com", "hash1")

		hash1, err := s.Hash()

		if err != nil {
			t.Fatalf("Error getting hash: %v", err)
		}

		if hash1 == "" {
			t.Fatalf("Expected hash to not be empty")
		}

		_, err = s.Create("page3", "http://example3.com", "hash1")

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

func TestUpdateQuietly(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo, nil)

	_, err := service.Create("1", "http://example.com", "hash1")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	page, err := service.Get("1")

	if err != nil {
		t.Fatalf("Error getting page: %v", err)
	}

	page.URL = "http://example2.com"

	err = service.UpdateQuietly(page)

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

func TestJSON(t *testing.T) {
	repo := pageRepo.NewRepository()

	s := service.NewService(repo, nil)
	_, err := s.Create("page1", "http://example.com", "hash1")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	_, err = s.Create("page2", "http://example2.com", "hash2")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	jsonB, err := s.JSON()

	if err != nil {
		t.Fatalf("Error getting json: %v", err)
	}

	var pages map[string]*domain.Page

	err = json.Unmarshal(jsonB, &pages)

	if err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("Expected 2 pages, got %d", len(pages))
	}

	if pages["page1"].ID != "page1" {
		t.Fatalf("Expected page ID to be 'page1', got '%s'", pages["page1"].ID)
	}

	if pages["page1"].URL != "http://example.com" {
		t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", pages["page1"].URL)
	}

}

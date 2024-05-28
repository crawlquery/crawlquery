package mem

import (
	"crawlquery/node/domain"
	"testing"
)

func TestPageRepo(t *testing.T) {
	r := NewRepository()
	err := r.Save("page1", &domain.Page{
		ID: "page1",
	})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	page, err := r.Get("page1")
	if err != nil {
		t.Fatalf("error getting page: %v", err)
	}

	if page == nil {
		t.Fatalf("expected page to be found")
	}

	if page.ID != "page1" {
		t.Fatalf("expected page ID to be 'page1', got '%s'", page.ID)
	}

	p2, err := r.Get("page2")

	if err == nil {
		t.Fatalf("expected error getting page2, got nil")
	}

	if p2 != nil {
		t.Fatalf("expected page2 to be nil, got %v", p2)
	}
}

func TestCount(t *testing.T) {
	r := NewRepository()
	r.Save("page1", &domain.Page{
		ID: "page1",
	})

	count, err := r.Count()
	if err != nil {
		t.Fatalf("error counting pages: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 page, got %d", count)
	}
}

func TestGetAll(t *testing.T) {
	r := NewRepository()
	r.Save("page1", &domain.Page{
		ID: "page1",
	})

	r.Save("page2", &domain.Page{
		ID: "page2",
	})

	pages, err := r.GetAll()
	if err != nil {
		t.Fatalf("error getting all pages: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}

	if pages["page1"].ID != "page1" {
		t.Fatalf("expected page1, got %s", pages["page1"].ID)
	}

	if pages["page2"].ID != "page2" {
		t.Fatalf("expected page2, got %s", pages["page2"].ID)
	}
}

func TestGetByIDs(t *testing.T) {
	r := NewRepository()
	r.Save("page1", &domain.Page{
		ID: "page1",
	})

	r.Save("page2", &domain.Page{
		ID: "page2",
	})

	r.Save("page3", &domain.Page{
		ID: "page3",
	})

	pages, err := r.GetByIDs([]string{"page1", "page2"})
	if err != nil {
		t.Fatalf("error getting pages by IDs: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}

	if pages["page1"].ID != "page1" {
		t.Fatalf("expected page1, got %s", pages["page1"].ID)
	}

	if pages["page2"].ID != "page2" {
		t.Fatalf("expected page2, got %s", pages["page2"].ID)
	}
}

func TestDelete(t *testing.T) {
	r := NewRepository()

	err := r.Save("page1", &domain.Page{
		ID: "page1",
	})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	err = r.Delete("page1")
	if err != nil {
		t.Fatalf("error deleting page: %v", err)
	}

	page, err := r.Get("page1")
	if err != domain.ErrPageNotFound {
		t.Fatalf("expected ErrPageNotFound, got %v", err)
	}

	if page != nil {
		t.Fatalf("expected page to be nil, got %v", page)
	}
}

func TestGetHashes(t *testing.T) {
	r := NewRepository()
	r.UpdateHash("page1", "hash1")

	hashes, err := r.GetHashes()
	if err != nil {
		t.Fatalf("error getting hashes: %v", err)
	}

	if len(hashes) != 1 {
		t.Fatalf("expected 1 hash, got %d", len(hashes))
	}

	if hashes["page1"] != "hash1" {
		t.Fatalf("expected hash1, got %s", hashes["page1"])
	}
}

func TestUpdateHash(t *testing.T) {
	r := NewRepository()
	err := r.UpdateHash("page1", "hash1")
	if err != nil {
		t.Fatalf("error updating hash: %v", err)
	}

	hash, err := r.GetHash("page1")
	if err != nil {
		t.Fatalf("error getting hash: %v", err)
	}

	if hash != "hash1" {
		t.Fatalf("expected hash1, got %s", hash)
	}
}

func TestDeletePageHash(t *testing.T) {
	r := NewRepository()
	r.UpdateHash("page1", "hash1")

	err := r.DeleteHash("page1")
	if err != nil {
		t.Fatalf("error deleting page: %v", err)
	}

	hash, err := r.GetHash("page1")
	if err != domain.ErrHashNotFound {
		t.Fatalf("expected ErrHashNotFound, got %v", err)
	}

	if hash != "" {
		t.Fatalf("expected empty string, got %s", hash)
	}
}

func TestGetHash(t *testing.T) {
	r := NewRepository()
	r.UpdateHash("page1", "hash1")

	hash, err := r.GetHash("page1")
	if err != nil {
		t.Fatalf("error getting hash: %v", err)
	}

	if hash != "hash1" {
		t.Fatalf("expected hash1, got %s", hash)
	}

	hash, err = r.GetHash("page2")
	if err != domain.ErrHashNotFound {
		t.Fatalf("expected ErrHashNotFound, got %v", err)
	}

	if hash != "" {
		t.Fatalf("expected empty string, got %s", hash)
	}
}

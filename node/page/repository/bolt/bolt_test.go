package bolt_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/page/repository/bolt"
	"os"
	"testing"
)

func TestRepo(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_get_test.db")
	defer os.Remove("/tmp/inverted_get_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.Save("page1", &domain.Page{
		ID:    "page1",
		URL:   "http://google.com",
		Title: "Google",
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
		t.Fatalf("expected page id to be page1, got %s", page.ID)
	}

	if page.URL != "http://google.com" {
		t.Fatalf("expected url to be http://google.com, got %s", page.URL)
	}

	if page.Title != "Google" {
		t.Fatalf("expected title to be Google, got %s", page.Title)
	}
}

func TestGetAll(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_get_all_test.db")
	defer os.Remove("/tmp/inverted_get_all_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.Save("page1", &domain.Page{
		ID:    "page1",
		URL:   "http://google.com",
		Title: "Google",
	})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	err = r.Save("page2", &domain.Page{
		ID:    "page2",
		URL:   "http://yahoo.com",
		Title: "Yahoo",
	})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

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

func TestDelete(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_remove_test.db")
	defer os.Remove("/tmp/inverted_remove_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.Save("page1", &domain.Page{
		ID:    "page1",
		URL:   "http://google.com",
		Title: "Google",
	})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	err = r.Delete("page1")

	if err != nil {
		t.Fatalf("error removing page: %v", err)
	}

	page, err := r.Get("page1")

	if err == nil {
		t.Fatalf("expected error getting page1, got nil")
	}

	if page != nil {
		t.Fatalf("expected page to be nil, got %v", page)
	}
}

func TestGetUpdateAndGetHash(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_get_update_hash_test.db")
	defer os.Remove("/tmp/inverted_get_update_hash_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.UpdateHash("page1", "hash1")

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

func TestGetHashes(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_get_hashes_test.db")
	defer os.Remove("/tmp/inverted_get_hashes_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.UpdateHash("page1", "hash1")

	if err != nil {
		t.Fatalf("error updating hash: %v", err)
	}

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

func TestDeleteHash(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_delete_hash_test.db")
	defer os.Remove("/tmp/inverted_delete_hash_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.UpdateHash("page1", "hash1")

	if err != nil {
		t.Fatalf("error updating hash: %v", err)
	}

	err = r.DeleteHash("page1")

	if err != nil {
		t.Fatalf("error deleting hash: %v", err)
	}

	_, err = r.GetHash("page1")

	if err == nil {
		t.Fatalf("expected error getting hash, got nil")
	}
}

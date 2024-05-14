package bolt_test

import (
	"crawlquery/node/index/forward/repository/bolt"
	"crawlquery/pkg/domain"
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

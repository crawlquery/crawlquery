package bolt_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/index/inverted/repository/bolt"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_get_test.db")
	defer os.Remove("/tmp/inverted_get_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.Save("keyword", &domain.Posting{PageID: "page1", Frequency: 1})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	page, err := r.Get("keyword")
	if err != nil {
		t.Fatalf("error getting page: %v", err)
	}

	if page == nil {
		t.Fatalf("expected page to be found")
	}

	if page[0].PageID != "page1" {
		t.Fatalf("expected page id to be page1, got %s", page[0].PageID)
	}

	if page[0].Frequency != 1 {
		t.Fatalf("expected frequency to be 1, got %d", page[0].Frequency)
	}
}

func TestSave(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_save_test.db")
	defer os.Remove("/tmp/inverted_save_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.Save("keyword", &domain.Posting{PageID: "page1", Frequency: 1})
	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}
}

func TestFuzzySearch(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_fuzzy_test.db")
	defer os.Remove("/tmp/inverted_fuzzy_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.Save("keyword", &domain.Posting{PageID: "page1", Frequency: 1})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	results := r.FuzzySearch("key")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results["page1"] != 1 {
		t.Fatalf("expected score to be 1, got %f", results["page1"])
	}
}

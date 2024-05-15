package bolt_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/keyword/repository/bolt"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_get_test.db")
	defer os.Remove("/tmp/inverted_get_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.SavePosting("keyword", &domain.Posting{PageID: "page1", Frequency: 1})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	postings, err := r.GetPostings("keyword")
	if err != nil {
		t.Fatalf("error getting page: %v", err)
	}

	if postings == nil {
		t.Fatalf("expected postings to be found")
	}

	if postings[0].PageID != "page1" {
		t.Fatalf("expected page id to be page1, got %s", postings[0].PageID)
	}

	if postings[0].Frequency != 1 {
		t.Fatalf("expected frequency to be 1, got %d", postings[0].Frequency)
	}
}

func TestSave(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_save_test.db")
	defer os.Remove("/tmp/inverted_save_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	err = r.SavePosting("keyword", &domain.Posting{PageID: "page1", Frequency: 1})
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

	err = r.SavePosting("keyword", &domain.Posting{PageID: "page1", Frequency: 1})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	results := r.FuzzySearch("key")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0] != "keyword" {
		t.Fatalf("expected result to be 'keyword', got '%s'", results[0])
	}
}

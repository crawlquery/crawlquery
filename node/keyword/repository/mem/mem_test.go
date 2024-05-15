package mem

import (
	"crawlquery/node/domain"
	"testing"
)

func TestSave(t *testing.T) {
	r := NewRepository()
	err := r.SavePosting("token", &domain.Posting{})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
}

func TestGet(t *testing.T) {
	r := NewRepository()
	err := r.SavePosting("token", &domain.Posting{})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	postings, err := r.GetPostings("token")
	if err != nil {
		t.Fatalf("error getting posting: %v", err)
	}
	if len(postings) != 1 {
		t.Fatalf("expected 1 posting, got %d", len(postings))
	}
}

func TestFuzzySearch(t *testing.T) {
	r := NewRepository()
	err := r.SavePosting("token", &domain.Posting{PageID: "page1", Frequency: 1})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	err = r.SavePosting("tom", &domain.Posting{PageID: "page2", Frequency: 2})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	results := r.FuzzySearch("to")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0] != "token" {
		t.Fatalf("expected first result to be 'token', got '%s'", results[0])
	}

	if results[1] != "tom" {
		t.Fatalf("expected second result to be 'tom', got '%s'", results[1])
	}
}

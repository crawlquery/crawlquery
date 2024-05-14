package mem

import (
	"crawlquery/node/domain"
	"testing"
)

func TestSave(t *testing.T) {
	r := NewRepository()
	err := r.Save("token", &domain.Posting{})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
}

func TestGet(t *testing.T) {
	r := NewRepository()
	err := r.Save("token", &domain.Posting{})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	postings, err := r.Get("token")
	if err != nil {
		t.Fatalf("error getting posting: %v", err)
	}
	if len(postings) != 1 {
		t.Fatalf("expected 1 posting, got %d", len(postings))
	}
}

func TestFuzzySearch(t *testing.T) {
	r := NewRepository()
	err := r.Save("token", &domain.Posting{PageID: "page1", Frequency: 1})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	err = r.Save("token", &domain.Posting{PageID: "page2", Frequency: 2})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	results := r.FuzzySearch("to")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results["page1"] != 1.0 {
		t.Fatalf("expected score of 1.0 for page1, got %f", results["page1"])
	}
	if results["page2"] != 2.0 {
		t.Fatalf("expected score of 2.0 for page2, got %f", results["page2"])
	}
}

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

func TestRemovePostingsByPageID(t *testing.T) {
	r := NewRepository()
	err := r.SavePosting("token", &domain.Posting{PageID: "page1", Frequency: 1})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	err = r.SavePosting("tom", &domain.Posting{PageID: "page2", Frequency: 2})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	err = r.SavePosting("token", &domain.Posting{PageID: "page3", Frequency: 3})
	if err != nil {
		t.Fatalf("error saving posting: %v", err)
	}
	err = r.RemovePostingsByPageID("page1")
	if err != nil {
		t.Fatalf("error removing postings: %v", err)
	}
	postings, err := r.GetPostings("token")
	if err != nil {
		t.Fatalf("error getting postings: %v", err)
	}
	if len(postings) != 1 {
		t.Fatalf("expected 1 posting, got %d", len(postings))
	}
	if postings[0].PageID != "page3" {
		t.Fatalf("expected posting to have pageID 'page3', got '%s'", postings[0].PageID)
	}
}

func TestUpdateHash(t *testing.T) {
	r := NewRepository()
	err := r.UpdateHash("token", "hash")
	if err != nil {
		t.Fatalf("error updating keyword hash: %v", err)
	}
	hash, err := r.GetHash("token")
	if err != nil {
		t.Fatalf("error getting keyword hash: %v", err)
	}
	if hash != "hash" {
		t.Fatalf("expected hash to be 'hash', got '%s'", hash)
	}
}

func TestGetHashes(t *testing.T) {
	r := NewRepository()
	err := r.UpdateHash("token", "hash")
	if err != nil {
		t.Fatalf("error updating keyword hash: %v", err)
	}
	hash, err := r.GetHash("token")
	if err != nil {
		t.Fatalf("error getting keyword hash: %v", err)
	}
	if hash != "hash" {
		t.Fatalf("expected hash to be 'hash', got '%s'", hash)
	}

	hashes, err := r.GetHashes()

	if err != nil {
		t.Fatalf("error getting keyword hashes: %v", err)
	}

	if len(hashes) != 1 {
		t.Fatalf("expected 1 hash, got %d", len(hashes))
	}
}

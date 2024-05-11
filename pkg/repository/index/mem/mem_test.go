package mem_test

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/index"
	"crawlquery/pkg/repository/index/mem"
	"reflect"
	"testing"
)

// hasForward checks if two indices have the same forward index.
func hasForward(idxA, idxB domain.Index) bool {
	if len(idxA.GetForward()) != len(idxB.GetForward()) {
		return false
	}
	for key, docA := range idxA.GetForward() {
		docB, ok := idxB.GetForward()[key]
		if !ok {
			return false
		}
		if docA.ID != docB.ID || docA.URL != docB.URL || docA.Title != docB.Title ||
			docA.Content != docB.Content || docA.MetaDescription != docB.MetaDescription {
			return false
		}
	}
	return true
}

// hasInverted checks if two indices have the same inverted index.
func hasInverted(idxA, idxB domain.Index) bool {
	if len(idxA.GetInverted()) != len(idxB.GetInverted()) {
		return false
	}
	for token, postingsA := range idxA.GetInverted() {
		postingsB, ok := idxB.GetInverted()[token]
		if !ok || len(postingsA) != len(postingsB) {
			return false
		}
		for i, postingA := range postingsA {
			postingB := postingsB[i]
			if postingA.PageID != postingB.PageID || postingA.Frequency != postingB.Frequency ||
				!reflect.DeepEqual(postingA.Positions, postingB.Positions) {
				return false
			}
		}
	}
	return true
}

func TestSaveAndLoadIndex(t *testing.T) {
	// Setup a test index and save to a temporary file
	testIndex := index.NewIndex()
	testIndex.AddPage(domain.Page{
		ID:      "doc1",
		URL:     "http://google.com",
		Title:   "Google",
		Content: "<html><body><h1>Google Search</h1></body></html>",
	})
	repo := mem.NewMemoryRepository()
	if err := repo.Save(testIndex); err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}

	// Attempt to load the index
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Failed to load index: %v", err)
	}

	if !hasForward(testIndex, loaded) {
		t.Fatalf("Forward index does not match expected")
	}

	if !hasInverted(testIndex, loaded) {
		t.Fatalf("Inverted index does not match expected")
	}

}

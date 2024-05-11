package disk_test

import (
	"crawlquery/pkg/index"
	"crawlquery/pkg/index/repository/disk"
	"os"
	"reflect"
	"testing"
)

// hasForward checks if two indices have the same forward index.
func hasForward(idxA, idxB *index.Index) bool {
	if len(idxA.Forward) != len(idxB.Forward) {
		return false
	}
	for key, docA := range idxA.Forward {
		docB, ok := idxB.Forward[key]
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
func hasInverted(idxA, idxB *index.Index) bool {
	if len(idxA.Inverted) != len(idxB.Inverted) {
		return false
	}
	for token, postingsA := range idxA.Inverted {
		postingsB, ok := idxB.Inverted[token]
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
	testIndex.AddDocument(index.Document{
		ID:      "doc1",
		URL:     "http://google.com",
		Title:   "Google",
		Content: "<html><body><h1>Google Search</h1></body></html>",
	})
	filepath := "/tmp/test_index.gob"
	repo := disk.NewDiskRepository(filepath)
	if err := repo.Save(testIndex); err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}

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

	// Cleanup
	os.ReadFile(filepath)
}

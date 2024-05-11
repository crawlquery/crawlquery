package index_test

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/index"
	"reflect"
	"sort"
	"testing"
)

func TestAddDocument(t *testing.T) {
	// Initialize index and document as before
	idx := index.NewIndex()
	doc := domain.Document{
		ID:              "doc1",
		URL:             "http://example.com",
		Title:           "Test Document",
		Content:         `<html><head><title>Example</title></head><body><h1>Hello World!</h1><p>This is a simple test. Numbers: 1234.</p></body></html>`,
		MetaDescription: "A simple test document",
	}

	// Add document to the index
	idx.AddDocument(doc)

	// Retrieve the document from the forward index
	indexedDoc, exists := idx.Forward[doc.ID]
	if !exists {
		t.Fatalf("Document with ID %s not found in forward index", doc.ID)
	}

	// Compare each field
	if indexedDoc.ID != doc.ID ||
		indexedDoc.URL != doc.URL ||
		indexedDoc.Title != doc.Title ||
		indexedDoc.Content != doc.Content ||
		indexedDoc.MetaDescription != doc.MetaDescription {
		t.Errorf("Document fields do not match. Got %+v, want %+v", indexedDoc, doc)
	}
}

func TestSearch(t *testing.T) {
	// Create a test index with some documents
	index := index.NewIndex()
	index.SetInverted(map[string][]domain.Posting{
		"test": {
			{PageID: "doc1", Frequency: 2},
			{PageID: "doc2", Frequency: 1},
		},
		"document": {
			{PageID: "doc1", Frequency: 1},
		},
	})

	// Define test cases
	tests := []struct {
		query      string
		wantResult []domain.Result
	}{
		{
			query: "test document",
			wantResult: []domain.Result{
				{PageID: "doc1", Score: 3}, // doc1 appears in both 'test' and 'document'
				{PageID: "doc2", Score: 1}, // doc2 appears only in 'test'
			},
		},
	}

	for _, tt := range tests {
		gotResults := index.Search(tt.query)

		// Need to sort results as slice order isn't guaranteed to match
		sort.Slice(gotResults, func(i, j int) bool {
			return gotResults[i].PageID < gotResults[j].PageID
		})

		if !reflect.DeepEqual(gotResults, tt.wantResult) {
			t.Errorf("Search(%q) = %v, want %v", tt.query, gotResults, tt.wantResult)
		}
	}
}

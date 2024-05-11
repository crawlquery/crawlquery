package index_test

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/index"
	"crawlquery/pkg/token"
	"reflect"
	"sort"
	"testing"
)

func TestAddPage(t *testing.T) {
	// Initialize index and page as before
	idx := index.NewIndex() // Assuming you have a constructor for Index
	doc := domain.Page{
		ID:              "doc1",
		URL:             "http://example.com",
		Title:           "Test Page",
		Content:         `<html><head><title>Example</title></head><body><h1>Hello World!</h1><p>This is a simple test. Numbers: 1234.</p></body></html>`,
		MetaDescription: "A simple test page",
	}

	// Tokenize the content to predict what should be in the inverted index
	// Assuming the existence of a Tokenize function that returns a map of token strings to their positions
	tokens := token.Tokenize(doc.Content) // Make sure to implement this or adjust to your actual tokenize function
	// Add page to the index
	idx.AddPage(doc)
	t.FailNow()
	// Retrieve the page from the forward index and verify it
	indexedDoc, exists := idx.Forward[doc.ID]
	if !exists {
		t.Fatalf("Page with ID %s not found in forward index", doc.ID)
	}
	if !reflect.DeepEqual(indexedDoc, doc) {
		t.Errorf("Page fields do not match. Got %+v, want %+v", indexedDoc, doc)
	}

	// Check the inverted index for correctness
	for token, positions := range tokens {
		postings, found := idx.Inverted[token]
		if !found {
			t.Errorf("Token %q not found in the inverted index", token)
			continue
		}

		// Check if the postings include the correct document with the correct frequency and positions
		var foundPosting bool
		for _, posting := range postings {
			if posting.PageID == doc.ID && reflect.DeepEqual(posting.Positions, positions) {
				foundPosting = true
				break
			}
		}
		if !foundPosting {
			t.Errorf("Posting for token %q with expected positions %v not found", token, positions)
		}
	}
}

func TestSearch(t *testing.T) {
	// Create a test index with some pages
	index := index.NewIndex()
	index.SetInverted(map[string][]domain.Posting{
		"test": {
			{PageID: "doc1", Frequency: 2},
			{PageID: "doc2", Frequency: 1},
		},
		"page": {
			{PageID: "doc1", Frequency: 1},
		},
	})

	index.SetForward(map[string]domain.Page{
		"doc1": {
			ID:              "doc1",
			URL:             "http://example.com/doc1",
			Title:           "Test Page",
			Content:         "This is a test page",
			MetaDescription: "A test page for indexing",
		},
		"doc2": {
			ID:              "doc2",
			URL:             "http://example.com/doc2",
			Title:           "Another Test",
			Content:         "This is another test",
			MetaDescription: "Another test page",
		},
	})

	// Define test cases
	tests := []struct {
		query      string
		wantResult []domain.Result
	}{
		{
			query: "test page",
			wantResult: []domain.Result{
				{PageID: "doc1", Score: 3, Page: index.Forward["doc1"]},
				{PageID: "doc2", Score: 1, Page: index.Forward["doc2"]},
			},
		},
	}

	for _, tt := range tests {
		gotResults, err := index.Search(tt.query)

		if err != nil {
			t.Fatalf("Search(%q) returned error: %v", tt.query, err)
		}

		// Need to sort results as slice order isn't guaranteed to match
		sort.Slice(gotResults, func(i, j int) bool {
			return gotResults[i].PageID < gotResults[j].PageID
		})

		if !reflect.DeepEqual(gotResults, tt.wantResult) {
			t.Errorf("Search(%q) = %v, want %v", tt.query, gotResults, tt.wantResult)
		}
	}
}

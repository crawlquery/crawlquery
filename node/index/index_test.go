package index

import (
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"crawlquery/node/token"
	sharedDomain "crawlquery/pkg/domain"
	"crawlquery/pkg/testutil"
	"fmt"

	keywordRepo "crawlquery/node/index/keyword/repository/mem"
	pageRepo "crawlquery/node/index/page/repository/mem"
	"reflect"
	"sort"
	"testing"
)

func TestAddPage(t *testing.T) {
	t.Run("should add a page to the index", func(t *testing.T) {
		// Initialize index and page as before
		idx := NewIndex(
			pageRepo.NewRepository(),
			keywordRepo.NewRepository(),
			testutil.NewTestLogger(),
		) // Assuming you have a constructor for Index
		doc := &sharedDomain.Page{
			ID:              "doc1",
			URL:             "http://example.com",
			Title:           "Test Page",
			MetaDescription: "A simple test page",
		}

		// Tokenize the content to predict what should be in the inverted index
		// Assuming the existence of a Tokenize function that returns a map of token strings to their positions
		tokens := token.Positions(doc.Keywords) // Make sure to implement this or adjust to your actual tokenize function
		// Add page to the index
		idx.AddPage(doc)
		// Retrieve the page from the forward index and verify it
		indexedDoc, err := idx.pageRepo.Get(doc.ID)
		if err != nil {
			t.Fatalf("Page with ID %s not found in forward index", doc.ID)
		}
		if !reflect.DeepEqual(indexedDoc, doc) {
			t.Errorf("Page fields do not match. Got %+v, want %+v", indexedDoc, doc)
		}

		// Check the inverted index for correctness
		for token, positions := range tokens {
			postings, err := idx.keywordRepo.Get(token)
			if err != nil {
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
	})

	t.Run("can add multiple pages to the index", func(t *testing.T) {
		// Initialize index and pages as before
		idx := NewIndex(
			pageRepo.NewRepository(),
			keywordRepo.NewRepository(),
			testutil.NewTestLogger(),
		) // Assuming you have a constructor for Index
		htmlData := [][]byte{
			[]byte("<html><body><p>Hello world!</p></body></html>"),
			[]byte("<html><body><p>Goodbye world!</p></body></html>"),
		}

		// Tokenize the content to predict what should be in the inverted index
		// Assuming the existence of a Tokenize function that returns a map of token strings to their positions
		for _, doc := range htmlData {

			page, err := parse.Parse("http://example.com", "doc1", doc)

			if err != nil {
				t.Fatalf("Error parsing HTML: %v", err)
			}

			tokens := token.Positions(page.Keywords) // Make sure to implement this or adjust to your actual tokenize function

			// Add page to the index
			idx.AddPage(page)
			// Retrieve the page from the forward index and verify it
			_, err := idx.pageRepo.Get(page.ID)
			if err != nil {
				t.Fatalf("Page with ID %s not found in pages repo", page.ID)
			}

			// Check the inverted index for correctness
			for token, positions := range tokens {
				postings, err := idx.keywordRepo.Get(token)
				if err != nil {
					t.Errorf("Token %q not found in the inverted index", token)
					continue
				}

				fmt.Println(postings)

				// Check if the postings include the correct document with the correct frequency and positions
				var foundPosting bool
				for _, posting := range postings {
					if posting.PageID == page.ID && reflect.DeepEqual(posting.Positions, positions) {
						foundPosting = true
						break
					}
				}
				if !foundPosting {
					t.Errorf("Posting for token %q with expected positions %v not found", token, positions)
				}
			}
		}
	})
}

func TestSearch(t *testing.T) {
	// Create a test index with some pages
	index := NewIndex(
		pageRepo.NewRepository(),
		keywordRepo.NewRepository(),
		testutil.NewTestLogger(),
	)
	inverted := map[string][]*domain.Posting{
		"test": {
			{PageID: "doc1", Frequency: 2},
			{PageID: "doc2", Frequency: 1},
		},
		"page": {
			{PageID: "doc1", Frequency: 1},
		},
	}

	forward := map[string]*sharedDomain.Page{
		"doc1": {
			ID:              "doc1",
			URL:             "http://example.com/doc1",
			Title:           "Test Page",
			MetaDescription: "A test page for indexing",
		},
		"doc2": {
			ID:              "doc2",
			URL:             "http://example.com/doc2",
			Title:           "Another Test",
			MetaDescription: "Another test page",
		},
	}

	for token, postings := range inverted {
		for _, posting := range postings {
			index.keywordRepo.Save(token, posting)
		}
	}

	for _, page := range forward {
		index.pageRepo.Save(page.ID, page)
	}

	// Define test cases
	tests := []struct {
		query      string
		wantResult []sharedDomain.Result
	}{
		{
			query: "test page",
			wantResult: []sharedDomain.Result{
				{PageID: "doc1", Score: 3, Page: forward["doc1"]},
				{PageID: "doc2", Score: 1, Page: forward["doc2"]},
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

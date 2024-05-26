package service

import (
	"reflect"
	"testing"

	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"
)

// Helper function to set up repositories and services
func setupTestRepos() (*pageRepo.Repository, *keywordRepo.Repository, *pageService.Service, *keywordService.Service) {
	keywordRepo := keywordRepo.NewRepository()
	keywordService := keywordService.NewService(keywordRepo)

	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo, nil)

	return pageRepo, keywordRepo, pageService, keywordService
}

// Helper function to save a page and its keywords
func savePage(t *testing.T, pageRepo *pageRepo.Repository, keywordRepo *keywordRepo.Repository, page domain.Page, keywords []string) {
	err := pageRepo.Save(page.ID, &page)
	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	err = keywordRepo.AddPageKeywords(page.ID, keywords)
	if err != nil {
		t.Fatalf("Error adding keywords: %v", err)
	}
}

func TestSearch(t *testing.T) {
	t.Run("searches for a single keyword", func(t *testing.T) {
		pageRepo, keywordRepo, pageService, keywordService := setupTestRepos()

		savePage(t, pageRepo, keywordRepo, domain.Page{
			ID:    "page1",
			URL:   "http://example.com",
			Title: "Example",
		}, []string{"example", "home page"})

		s := NewService(pageService, keywordService)
		results, err := s.Search("example")

		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].PageID != "page1" {
			t.Errorf("Expected page1, got %s", results[0].PageID)
		}
	})

	t.Run("searches for multiple keywords", func(t *testing.T) {
		pageRepo, keywordRepo, pageService, keywordService := setupTestRepos()

		pages := []struct {
			id, url, title string
			keywords       []string
		}{
			{"page1", "http://example.com", "Example", []string{"example", "home page"}},
			{"page2", "http://example.com/contact", "Contact", []string{"example", "contact"}},
		}

		for _, p := range pages {
			savePage(t, pageRepo, keywordRepo, domain.Page{
				ID:    p.id,
				URL:   p.url,
				Title: p.title,
			}, p.keywords)
		}

		s := NewService(pageService, keywordService)
		results, err := s.Search("example contact")

		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("Expected 2 results, got %d", len(results))
		}

		// page2 has more hits
		expectedIDs := []string{"page2", "page1"}
		for i, expectedID := range expectedIDs {
			if results[i].PageID != expectedID {
				t.Errorf("Expected %s, got %s", expectedID, results[i].PageID)
			}
		}
	})

	t.Run("searches for long keyword", func(t *testing.T) {
		pageRepo, keywordRepo, pageService, keywordService := setupTestRepos()

		savePage(t, pageRepo, keywordRepo, domain.Page{
			ID:    "page1",
			URL:   "http://example.com",
			Title: "Example",
		}, []string{"best way to detect bot from user agent"})

		s := NewService(pageService, keywordService)
		results, err := s.Search("best way to detect bot from user agent")

		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].PageID != "page1" {
			t.Errorf("Expected page1, got %s", results[0].PageID)
		}
	})

	t.Run("searches for two separate keywords", func(t *testing.T) {
		pageRepo, keywordRepo, pageService, keywordService := setupTestRepos()

		savePage(t, pageRepo, keywordRepo, domain.Page{
			ID:    "page1",
			URL:   "http://example.com",
			Title: "Example",
		}, []string{"market"})

		savePage(t, pageRepo, keywordRepo, domain.Page{
			ID:    "page2",
			URL:   "http://example.com",
			Title: "Example",
		}, []string{"data"})

		s := NewService(pageService, keywordService)
		results, err := s.Search("market data")

		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("Expected 2 results, got %d", len(results))
		}

		expectedIDs := []string{"page1", "page2"}
		for i, expectedID := range expectedIDs {
			if results[i].PageID != expectedID {
				t.Errorf("Expected %s, got %s", expectedID, results[i].PageID)
			}
		}
	})
}

func TestGetResultsForTermGroup(t *testing.T) {
	pageRepo, keywordRepo, pageService, keywordService := setupTestRepos()

	savePage(t, pageRepo, keywordRepo, domain.Page{
		ID:    "page1",
		URL:   "http://example.com",
		Title: "Example",
	}, []string{"example", "home page"})

	savePage(t, pageRepo, keywordRepo, domain.Page{
		ID:    "page2",
		URL:   "http://example.com/contact",
		Title: "Contact",
	}, []string{"example", "contact"})

	s := NewService(pageService, keywordService)
	termGroups := [][]string{{"example"}, {"contact"}}
	results, err := s.getResultsForTermGroup(termGroups)

	if err != nil {
		t.Fatalf("Error getting results: %v", err)
	}

	expectedHits := map[string]map[string]int{
		"page1": {"example": 1},
		"page2": {"example": 1, "contact": 1},
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	for _, result := range results {
		expectedHitMap, ok := expectedHits[result.PageID]
		if !ok {
			t.Fatalf("Unexpected page ID %s in results", result.PageID)
		}
		if !reflect.DeepEqual(result.Hits, expectedHitMap) {
			t.Errorf("Expected hits %v for page %s, got %v", expectedHitMap, result.PageID, result.Hits)
		}
	}
}

func TestSplitQueryIntoCombinations(t *testing.T) {
	tests := []struct {
		query    string
		expected [][]string
	}{
		{
			query: "search service",
			expected: [][]string{
				{"search"},
				{"search", "service"},
				{"service"},
			},
		},
		{
			query: "split query into combinations",
			expected: [][]string{
				{"split"},
				{"split", "query"},
				{"split", "query", "into"},
				{"split", "query", "into", "combinations"},
				{"query"},
				{"query", "into"},
				{"query", "into", "combinations"},
				{"into"},
				{"into", "combinations"},
				{"combinations"},
			},
		},
		{
			query: "oneword",
			expected: [][]string{
				{"oneword"},
			},
		},
		{
			query:    "",
			expected: [][]string{},
		},
	}

	for _, test := range tests {
		result := splitQueryIntoCombinations(test.query)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For query '%s', expected %v, but got %v", test.query, test.expected, result)
		}
	}
}

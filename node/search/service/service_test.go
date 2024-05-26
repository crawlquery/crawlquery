package service_test

import (
	"reflect"
	"testing"

	"crawlquery/node/domain"
	"crawlquery/node/html/repository/mem"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/occurrence/service"
)

func setupTestRepos() (*pageRepo.Repository, *mem.Repository, *pageService.Service, *service.KeywordService) {
	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo, nil)

	keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
	keywordService := keywordService.NewKeywordService(keywordRepo)

	return pageRepo, keywordRepo, pageService, keywordService
}

func savePage(t *testing.T, pageRepo *pageRepo.Repository, keywordRepo *mem.Repository, page domain.Page, keywordOccurrences map[domain.Keyword]domain.KeywordOccurrence) {
	err := pageRepo.Save(page.ID, &page)
	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	for keyword, occurrence := range keywordOccurrences {
		err := keywordRepo.Add(keyword, occurrence)
		if err != nil {
			t.Fatalf("Error adding keyword occurrence: %v", err)
		}
	}
}

func TestService_Search(t *testing.T) {
	pageRepo, keywordRepo, pageService, keywordService := setupTestRepos()
	svc := service.NewService(pageService, keywordService)

	// Add pages and keyword occurrences
	page1 := domain.Page{ID: "page1", URL: "http://example.com", Title: "Example"}
	page2 := domain.Page{ID: "page2", URL: "http://example.com/contact", Title: "Contact"}

	keywordOccurrences1 := map[domain.Keyword]domain.KeywordOccurrence{
		"example": {PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
	}
	keywordOccurrences2 := map[domain.Keyword]domain.KeywordOccurrence{
		"example": {PageID: "page2", Frequency: 1, Positions: []int{1}},
		"contact": {PageID: "page2", Frequency: 1, Positions: []int{4}},
	}

	savePage(t, pageRepo, keywordRepo, page1, keywordOccurrences1)
	savePage(t, pageRepo, keywordRepo, page2, keywordOccurrences2)

	// Test Search
	results, err := svc.Search("example contact")
	if err != nil {
		t.Fatalf("Error searching: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	expectedResults := []*domain.Result{
		{
			PageID: "page1",
			Page: &domain.ResultPage{
				ID:    "page1",
				Hash:  "",
				URL:   "http://example.com",
				Title: "Example",
			},
			Score: 3,
			KeywordOccurences: map[string]domain.KeywordOccurrence{
				"example": {PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
			},
		},
		{
			PageID: "page2",
			Page: &domain.ResultPage{
				ID:    "page2",
				Hash:  "",
				URL:   "http://example.com/contact",
				Title: "Contact",
			},
			Score: 2,
			KeywordOccurences: map[string]domain.KeywordOccurrence{
				"example": {PageID: "page2", Frequency: 1, Positions: []int{1}},
				"contact": {PageID: "page2", Frequency: 1, Positions: []int{4}},
			},
		},
	}

	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}
}

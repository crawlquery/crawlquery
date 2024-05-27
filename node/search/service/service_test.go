package service_test

import (
	"reflect"
	"testing"

	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"
	"crawlquery/node/search/service"
	"crawlquery/pkg/testutil"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"
)

func setupTestRepos() (*pageRepo.Repository, *keywordOccurrenceRepo.Repository, *pageService.Service, *keywordService.Service) {
	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo, nil)

	keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
	keywordService := keywordService.NewService(keywordOccurrenceRepo)

	return pageRepo, keywordOccurrenceRepo, pageService, keywordService
}

func savePage(t *testing.T, pageRepo *pageRepo.Repository, keywordOccurrenceRepo *keywordOccurrenceRepo.Repository, page domain.Page, keywordOccurrences map[domain.Keyword]domain.KeywordOccurrence) {
	err := pageRepo.Save(page.ID, &page)
	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	for keyword, occurrence := range keywordOccurrences {
		err := keywordOccurrenceRepo.Add(keyword, occurrence)
		if err != nil {
			t.Fatalf("Error adding keyword occurrence: %v", err)
		}
	}
}

func checkResult(t *testing.T, result, expected domain.Result) {
	if result.PageID != expected.PageID {
		t.Errorf("Expected PageID %s, got %s", expected.PageID, result.PageID)
	}
	if result.Page.ID != expected.Page.ID {
		t.Errorf("Expected Page ID %s, got %s", expected.Page.ID, result.Page.ID)
	}
	if result.Page.URL != expected.Page.URL {
		t.Errorf("Expected Page URL %s, got %s", expected.Page.URL, result.Page.URL)
	}
	if result.Page.Title != expected.Page.Title {
		t.Errorf("Expected Page Title %s, got %s", expected.Page.Title, result.Page.Title)
	}
	if result.Score != expected.Score {
		t.Errorf("Expected Score %f, got %f", expected.Score, result.Score)
	}
	if !reflect.DeepEqual(result.KeywordOccurences, expected.KeywordOccurences) {
		t.Errorf("Expected KeywordOccurences %v, got %v", expected.KeywordOccurences, result.KeywordOccurences)
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

	expectedResults := []domain.Result{
		{
			PageID: "page2",
			Page: domain.ResultPage{
				ID:    "page2",
				Hash:  "", // Assuming Hash is not being set for the test
				URL:   "http://example.com/contact",
				Title: "Contact",
			},
			Score: 4,
			KeywordOccurences: map[string]domain.KeywordOccurrence{
				"example": {PageID: "page2", Frequency: 1, Positions: []int{1}},
				"contact": {PageID: "page2", Frequency: 1, Positions: []int{4}},
			},
		},
		{
			PageID: "page1",
			Page: domain.ResultPage{
				ID:    "page1",
				Hash:  "", // Assuming Hash is not being set for the test
				URL:   "http://example.com",
				Title: "Example",
			},
			Score: 3,
			KeywordOccurences: map[string]domain.KeywordOccurrence{
				"example": {PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
			},
		},
	}

	testutil.PrettyPrint(results)
	for i, result := range results {
		checkResult(t, result, expectedResults[i])
	}
}

package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"crawlquery/node/domain"
	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"
	searchHandler "crawlquery/node/search/handler"
	searchService "crawlquery/node/search/service"
	"crawlquery/pkg/testutil"
)

func setupTestRepos() (*pageRepo.Repository, *keywordOccurrenceRepo.Repository, *pageService.Service, *keywordService.Service) {
	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo, nil)

	keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
	keywordService := keywordService.NewService(keywordOccurrenceRepo)

	return pageRepo, keywordOccurrenceRepo, pageService, keywordService
}

func savePage(t *testing.T, pageRepo *pageRepo.Repository, keywordRepo *keywordOccurrenceRepo.Repository, page domain.Page, keywordOccurrences map[domain.Keyword]domain.KeywordOccurrence) {
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

func checkResponseBody(t *testing.T, body string, expectedPageIDs []string) {
	for _, pageID := range expectedPageIDs {
		if !strings.Contains(body, pageID) {
			t.Errorf("expected body to contain '%s'; got %s", pageID, body)
		}
	}
}

func TestSearch(t *testing.T) {

	t.Run("returns results", func(t *testing.T) {
		pageRepo, keywordRepo, pageService, keywordService := setupTestRepos()
		searchSvc := searchService.NewService(pageService, keywordService)

		page1 := domain.Page{
			ID:    "page1",
			URL:   "http://example.com",
			Title: "Example",
		}
		page2 := domain.Page{
			ID:    "page2",
			URL:   "http://example.com/contact",
			Title: "Contact",
		}

		keywordOccurrences1 := map[domain.Keyword]domain.KeywordOccurrence{
			"keyword": {PageID: "page1", Frequency: 1, Positions: []int{1}},
		}
		keywordOccurrences2 := map[domain.Keyword]domain.KeywordOccurrence{
			"keyword": {PageID: "page2", Frequency: 1, Positions: []int{1}},
		}

		savePage(t, pageRepo, keywordRepo, page1, keywordOccurrences1)
		savePage(t, pageRepo, keywordRepo, page2, keywordOccurrences2)

		searchHandler := searchHandler.NewHandler(searchSvc, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search?q=keyword", nil)

		searchHandler.Search(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		checkResponseBody(t, w.Body.String(), []string{"page1", "page2"})
	})

	t.Run("returns error if query is missing", func(t *testing.T) {
		searchHandler := searchHandler.NewHandler(nil, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search", nil)

		searchHandler.Search(ctx)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status BadRequest; got %v", w.Code)
		}
	})
}

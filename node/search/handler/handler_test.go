package handler_test

import (
	"crawlquery/pkg/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	searchService "crawlquery/node/search/service"

	searchHandler "crawlquery/node/search/handler"
)

func TestSearch(t *testing.T) {

	t.Run("returns results", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		searchService := searchService.NewService(pageService, keywordService)

		err := pageRepo.Save("page1", &domain.Page{
			ID:    "page1",
			URL:   "http://example.com",
			Title: "Example",
		})

		if err != nil {
			t.Errorf("Error saving page: %v", err)
		}

		err = keywordRepo.AddPageKeywords("page1", []string{"market"})
		if err != nil {
			t.Errorf("Error adding keywords: %v", err)
		}

		err = pageRepo.Save("page2", &domain.Page{
			ID:    "page2",
			URL:   "http://example.com",
			Title: "Example",
		})

		if err != nil {
			t.Errorf("Error saving page: %v", err)
		}

		err = keywordRepo.AddPageKeywords("page2", []string{"data"})

		if err != nil {
			t.Errorf("Error adding keywords: %v", err)
		}

		indexHandler := searchHandler.NewHandler(searchService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search?q=market data", nil)

		indexHandler.Search(ctx)

		body := w.Body.String()

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		if !strings.Contains(body, "page1") {
			t.Errorf("expected body to contain 'home1'; got %s", body)
		}
	})

	t.Run("returns error if query is missing", func(t *testing.T) {

		indexHandler := searchHandler.NewHandler(nil, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search", nil)

		indexHandler.Search(ctx)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status BadRequest; got %v", w.Code)
		}
	})
}

package handler_test

import (
	indexHandler "crawlquery/node/index/handler"
	"crawlquery/pkg/domain"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	indexService "crawlquery/node/index/service"

	"github.com/gin-gonic/gin"
)

func TestSearch(t *testing.T) {

	t.Run("returns results", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		indexService := indexService.NewService(
			pageService,
			htmlService,
			keywordService,
			testutil.NewTestLogger(),
		)

		for k, dummy := range factory.ThreePages {

			page, err := pageService.Create(k, dummy.URL)

			if err != nil {
				t.Fatalf("error saving page: %v", err)
			}

			err = htmlService.Save(page.ID, []byte(dummy.HTML))

			if err != nil {
				t.Fatalf("error saving html: %v", err)
			}

			err = indexService.Index(page.ID)

			if err != nil {
				t.Fatalf("error indexing page: %v", err)
			}
		}

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search?q=example+website", nil)

		indexHandler.Search(ctx)

		expected, err := json.Marshal(gin.H{
			"results": []domain.Result{
				{
					PageID: "home1",
					Score:  1.0,
					Page: &domain.Page{
						ID:              "home1",
						URL:             "https://example.com",
						Title:           "Home",
						MetaDescription: "Welcome to our official website where we offer the latest updates and information.",
					},
				},
			},
		})

		if err != nil {
			t.Fatal(err)
		}

		body := w.Body.String()

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		if body != string(expected) {
			t.Errorf("expected body %s; got %s", expected, body)
		}
	})

	t.Run("returns error if query is missing", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		indexService := indexService.NewService(
			pageService,
			htmlService,
			keywordService,
			testutil.NewTestLogger(),
		)

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search", nil)

		indexHandler.Search(ctx)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status BadRequest; got %v", w.Code)
		}
	})
}

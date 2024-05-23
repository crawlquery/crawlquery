package handler_test

import (
	"bytes"
	"crawlquery/node/domain"
	indexHandler "crawlquery/node/index/handler"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	indexService "crawlquery/node/index/service"

	peerService "crawlquery/node/peer/service"

	"github.com/gin-gonic/gin"
)

func TestSearch(t *testing.T) {

	t.Run("returns results", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		peerService := peerService.NewService(nil, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
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

		body := w.Body.String()

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		if !strings.Contains(body, "home1") {
			t.Errorf("expected body to contain 'home1'; got %s", body)
		}
	})

	t.Run("returns error if query is missing", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		peerService := peerService.NewService(nil, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
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

func TestReIndex(t *testing.T) {
	t.Run("reindexes page", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		peerService := peerService.NewService(nil, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
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
		}

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/reindex/home1", nil)
		ctx.Params = gin.Params{
			{Key: "pageID", Value: "home1"},
		}

		indexHandler.ReIndex(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		page, err := pageService.Get("home1")

		if err != nil {
			t.Fatalf("error getting page: %v", err)
		}

		if page == nil {
			t.Fatalf("expected page to be found, got nil")
		}

		if page.ID != "home1" {
			t.Fatalf("expected page ID to be home1, got %s", page.ID)
		}

		if page.Title != "Home" {
			t.Fatalf("expected title to be Home, got %s", page.Title)
		}
	})
}

func TestGetIndex(t *testing.T) {
	t.Run("returns page", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		peerService := peerService.NewService(nil, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
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

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/index/home1", nil)
		ctx.Params = gin.Params{
			{Key: "pageID", Value: "home1"},
		}

		indexHandler.GetIndex(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		var response *domain.Page

		err := json.Unmarshal(w.Body.Bytes(), &response)

		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		if response.ID != "home1" {
			t.Fatalf("expected page ID to be home1, got %s", response.ID)
		}

		if response.Title != "Home" {
			t.Fatalf("expected title to be Home, got %s", response.Title)
		}
	})
}

func TestEvent(t *testing.T) {
	t.Run("can handle index event", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		peerService := peerService.NewService(nil, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
			testutil.NewTestLogger(),
		)

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		indexEvent := domain.IndexEvent{
			Page: &domain.Page{
				URL:             "http://example.com",
				ID:              "page1",
				Title:           "Example",
				MetaDescription: "An example page",
				Keywords:        []string{"distro", "linux"},
			},
		}

		encoded, err := json.Marshal(indexEvent)

		if err != nil {
			t.Fatalf("error encoding index event: %v", err)
		}

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(encoded))

		indexHandler.Event(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		page, err := pageService.Get("page1")

		if err != nil {
			t.Fatalf("error getting page: %v", err)
		}

		if page == nil {
			t.Fatalf("expected page to be found, got nil")
		}

		if page.ID != "page1" {
			t.Fatalf("expected page ID to be page1, got %s", page.ID)
		}

		if page.Title != "Example" {
			t.Fatalf("expected title to be Example, got %s", page.Title)
		}

		if page.MetaDescription != "An example page" {
			t.Fatalf("expected meta description to be An example page, got %s", page.MetaDescription)
		}

		if len(page.Keywords) != 2 {
			t.Fatalf("expected 2 keywords, got %d", len(page.Keywords))
		}

		if page.Keywords[0] != "distro" {
			t.Fatalf("expected keyword to be distro, got %s", page.Keywords[0])
		}

		if page.Keywords[1] != "linux" {
			t.Fatalf("expected keyword to be linux, got %s", page.Keywords[1])
		}
	})
}

func TestHash(t *testing.T) {
	t.Run("returns hash", func(t *testing.T) {

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		indexService := indexService.NewService(pageService, nil, nil, testutil.NewTestLogger())

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/hash", nil)

		indexHandler.Hash(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		var response map[string]string

		err := json.Unmarshal(w.Body.Bytes(), &response)

		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		expectedPageHash, err := pageService.Hash()

		if err != nil {
			t.Fatalf("error getting page hash: %v", err)
		}

		if response["page"] != expectedPageHash {
			t.Fatalf("expected hash %s; got %s", expectedPageHash, response["page"])
		}

	})
}

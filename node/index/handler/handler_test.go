package handler_test

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/dto"
	indexHandler "crawlquery/node/index/handler"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	backupService "crawlquery/node/html/backup/service"
	htmlClient "crawlquery/pkg/client/html"

	indexService "crawlquery/node/index/service"

	peerService "crawlquery/node/peer/service"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
)

func setupTestRepos() (
	*pageRepo.Repository,
	*pageService.Service,
	*htmlRepo.Repository,
	*htmlService.Service,
	*peerService.Service,
	*keywordService.Service,
) {
	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo, nil)

	backupService := backupService.NewService(htmlClient.NewClient("http://storage:8080"))

	htmlRepo := htmlRepo.NewRepository()
	htmlService := htmlService.NewService(htmlRepo, backupService)

	peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

	keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
	keywordService := keywordService.NewService(keywordOccurrenceRepo)

	return pageRepo, pageService, htmlRepo, htmlService, peerService, keywordService
}

func TestIndex(t *testing.T) {
	t.Run("indexes page", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://storage:8080").
			Post("/pages/hash1").
			Reply(201)

		pageRepo, pageService, htmlRepo, htmlService, peerService, keywordService := setupTestRepos()

		pageRepo.Save("home1", &domain.Page{
			ID:  "home1",
			URL: "http://example.com",
		})

		htmlRepo.Save("hash1", []byte(`<html><head><title>Home</title></head><body><p>Hello this is my home page</p></html>`))

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
			keywordService,
			testutil.NewTestLogger(),
		)

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := dto.IndexRequest{
			PageID:      "home1",
			URL:         "http://example.com",
			ContentHash: "hash1",
		}

		encoded, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("error encoding index request: %v", err)
		}

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/index", bytes.NewBuffer(encoded))

		indexHandler.Index(ctx)

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

		_, pageService, htmlRepo, htmlService, peerService, keywordService := setupTestRepos()

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
			keywordService,
			testutil.NewTestLogger(),
		)

		for k, dummy := range factory.ThreePages {

			page, err := pageService.Create(k, dummy.URL, "hash1")

			if err != nil {
				t.Fatalf("error saving page: %v", err)
			}

			err = htmlRepo.Save(page.Hash, []byte(dummy.HTML))

			if err != nil {
				t.Fatalf("error saving html: %v", err)
			}

			err = indexService.Index(page.ID, page.URL, page.Hash)

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

		_, pageService, _, htmlService, peerService, keywordService := setupTestRepos()

		indexService := indexService.NewService(
			pageService,
			htmlService,
			peerService,
			keywordService,
			testutil.NewTestLogger(),
		)

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		PageUpdatedEvent := domain.PageUpdatedEvent{
			Page: &domain.Page{
				URL:         "http://example.com",
				ID:          "page1",
				Title:       "Example",
				Description: "An example page",
			},
		}

		encoded, err := json.Marshal(PageUpdatedEvent)

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

		if page.Description != "An example page" {
			t.Fatalf("expected meta description to be An example page, got %s", page.Description)
		}

	})
}

func TestHash(t *testing.T) {
	t.Run("returns hash", func(t *testing.T) {

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		indexService := indexService.NewService(pageService, nil, nil, nil, testutil.NewTestLogger())

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

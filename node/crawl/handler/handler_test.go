package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	crawlHandler "crawlquery/node/crawl/handler"
	crawlService "crawlquery/node/crawl/service"
	"crawlquery/node/dto"
	htmlBackupService "crawlquery/node/html/backup/service"
	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"
	indexService "crawlquery/node/index/service"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"
	peerService "crawlquery/node/peer/service"
	"crawlquery/pkg/client/api"
	"crawlquery/pkg/client/html"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
)

func setupServices() (*crawlService.CrawlService, *htmlRepo.Repository, *pageRepo.Repository) {
	htmlRepository := htmlRepo.NewRepository()
	htmlBackupSvc := htmlBackupService.NewService(html.NewClient("http://storage:8080"))
	htmlSvc := htmlService.NewService(htmlRepository, htmlBackupSvc)

	pageRepository := pageRepo.NewRepository()
	pageSvc := pageService.NewService(pageRepository, nil)

	peerSvc := peerService.NewService(nil, nil, testutil.NewTestLogger())
	indexSvc := indexService.NewService(pageSvc, htmlSvc, peerSvc, nil, testutil.NewTestLogger())
	apiClient := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

	crawlSvc := crawlService.NewService(htmlSvc, pageSvc, indexSvc, apiClient, testutil.NewTestLogger())

	return crawlSvc, htmlRepository, pageRepository
}

func TestCrawl(t *testing.T) {
	t.Run("can crawl a page", func(t *testing.T) {
		defer gock.Off()
		crawlSvc, htmlRepo, pageRepo := setupServices()

		expectedData := "<html><head><title>Example</title></head><body><h1>Hello, World!</h1><p>Website description</p></body></html>"
		expectedPageHash := util.Sha256Hex32([]byte(expectedData))

		gock.New("http://storage:8080").
			Post("/pages").
			Reply(201)

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			Header.Set("Content-Type", "text/html")

		req := dto.CrawlRequest{
			PageID: "test1",
			URL:    "http://example.com",
		}

		body, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Error marshalling request: %v", err)
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/crawl", bytes.NewBuffer(body))

		handler := crawlHandler.NewHandler(crawlSvc, testutil.NewTestLogger())
		handler.Crawl(ctx)

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
		}

		var resp dto.CrawlResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if resp.Page.Hash != expectedPageHash {
			t.Fatalf("Expected page hash to be '%s', got '%s'", expectedPageHash, resp.Page.Hash)
		}

		data, err := htmlRepo.Get("test1")
		if err != nil {
			t.Fatalf("Error reading data: %v", err)
		}

		if string(data) != expectedData {
			t.Fatalf("Expected data to be '%s', got '%s'", expectedData, data)
		}

		page, err := pageRepo.Get("test1")
		if err != nil {
			t.Fatalf("Error getting page: %v", err)
		}

		if page == nil {
			t.Fatalf("Expected page to be found")
		}

		if page.ID != "test1" {
			t.Fatalf("Expected page ID to be 'test1', got '%s'", page.ID)
		}

		if page.URL != "http://example.com" {
			t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", page.URL)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})

	t.Run("handles 404", func(t *testing.T) {
		defer gock.Off()
		crawlSvc, _, pageRepo := setupServices()

		gock.New("http://example.com").
			Get("/").
			Reply(404)

		req := dto.CrawlRequest{
			PageID: "test1",
			URL:    "http://example.com",
		}

		body, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Error marshalling request: %v", err)
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/crawl", bytes.NewBuffer(body))

		handler := crawlHandler.NewHandler(crawlSvc, testutil.NewTestLogger())
		handler.Crawl(ctx)

		if w.Code != 400 {
			t.Errorf("Expected status 400, got %v", w.Code)
		}

		page, err := pageRepo.Get("test1")
		if err == nil || page != nil {
			t.Fatalf("Expected page to be nil")
		}
	})
}

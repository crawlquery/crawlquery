package handler_test

import (
	"bytes"
	"crawlquery/node/dto"
	"crawlquery/pkg/client/html"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"encoding/json"
	"net/http/httptest"
	"testing"

	crawlHandler "crawlquery/node/crawl/handler"

	htmlBackupService "crawlquery/node/html/backup/service"
	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	peerService "crawlquery/node/peer/service"

	indexService "crawlquery/node/index/service"

	crawlService "crawlquery/node/crawl/service"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	t.Run("can crawl a page", func(t *testing.T) {
		defer gock.Off()
		htmlRepo := htmlRepo.NewRepository()

		gock.New("http://storage:8080").
			Post("/pages").
			Reply(201)

		htmlBackupService := htmlBackupService.NewService(html.NewClient("http://storage:8080"))

		htmlService := htmlService.NewService(htmlRepo, htmlBackupService)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		indexService := indexService.NewService(pageService, htmlService, peerService, keywordService, testutil.NewTestLogger())

		crawlService := crawlService.NewService(htmlService, pageService, indexService, nil, testutil.NewTestLogger())

		expectedData := "<html><head><title>Example</title></head><body><h1>Hello, World!</h1><p>Website description</p></body></html>"

		expectedPageHash := util.Sha256Hex32([]byte(expectedData))

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

		handler := crawlHandler.NewHandler(crawlService, testutil.NewTestLogger())

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
		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo, nil)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		crawlService := crawlService.NewService(htmlService, pageService, nil, nil, testutil.NewTestLogger())

		defer gock.Off()

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

		handler := crawlHandler.NewHandler(crawlService, testutil.NewTestLogger())

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

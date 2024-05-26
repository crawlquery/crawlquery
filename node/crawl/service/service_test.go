package service_test

import (
	"crawlquery/node/crawl/service"
	htmlBackupService "crawlquery/node/html/backup/service"
	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	peerService "crawlquery/node/peer/service"

	indexService "crawlquery/node/index/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	"crawlquery/pkg/client/api"
	"crawlquery/pkg/client/html"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"

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

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(htmlService, pageService, indexService, api, testutil.NewTestLogger())

		expectedData := "<html><head><title>Example</title></head><body><h1>Hello, World!</h1><p>Welcome to my example website.</p></body></html>"

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).Header.Set("Content-Type", "text/html")

		_, err := service.Crawl("test1", "http://example.com")

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
		}

		if err != nil {
			t.Fatalf("Error creating repository: %v", err)
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
			t.Errorf("Expected all mocks to be called")
		}
	})

	t.Run("creates crawl jobs for links", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://storage:8080").
			Post("/pages").
			Reply(201)

		htmlRepo := htmlRepo.NewRepository()
		htmlBackupService := htmlBackupService.NewService(html.NewClient("http://storage:8080"))
		htmlService := htmlService.NewService(htmlRepo, htmlBackupService)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		indexService := indexService.NewService(pageService, htmlService, peerService, keywordService, testutil.NewTestLogger())

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(htmlService, pageService, indexService, api, testutil.NewTestLogger())

		expectedData := `<html><head><title>Example</title></head><body><h1>Hello, World! <a href="http://example.com/about">About us</a></h1><p>Welcome to my website.</p></body></html>`

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			SetHeader("Content-Type", "text/html")

		gock.New("http://localhost:8080").
			Post("/links").
			JSON(`{"src":"http://example.com","dst":"http://example.com/about"}`).
			Reply(200)

		pageCrawled, err := service.Crawl("test1", "http://example.com")

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
		}

		if pageCrawled == nil {
			t.Fatalf("Expected page to be crawled")
		}

		expectedHash := util.Sha256Hex32([]byte(expectedData))

		if pageCrawled.Hash != expectedHash {
			t.Fatalf("Expected page hash to be %s, got %s", expectedHash, pageCrawled.Hash)
		}

		if err != nil {
			t.Fatalf("Error creating repository: %v", err)
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
			t.Errorf("Expected all mocks to be called")
		}
	})

	t.Run("handles relative links", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://storage:8080").
			Post("/pages").
			Reply(201)

		htmlRepo := htmlRepo.NewRepository()
		htmlBackupService := htmlBackupService.NewService(html.NewClient("http://storage:8080"))
		htmlService := htmlService.NewService(htmlRepo, htmlBackupService)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		indexService := indexService.NewService(pageService, htmlService, peerService, keywordService, testutil.NewTestLogger())

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(htmlService, pageService, indexService, api, testutil.NewTestLogger())

		expectedData := `<html><head><title>Example</title></head><body><h1>Hello, World! <a href="/about">About us</a></h1><p>Welcome to my website.</p></body></html>`

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			SetHeader("Content-Type", "text/html")

		gock.New("http://localhost:8080").
			Post("/links").
			JSON(`{"src":"http://example.com","dst":"http://example.com/about"}`).
			Reply(200)

		pageCrawled, err := service.Crawl("test1", "http://example.com")

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
		}

		if pageCrawled == nil {
			t.Fatalf("Expected page to be crawled")
		}

		expectedHash := util.Sha256Hex32([]byte(expectedData))

		if pageCrawled.Hash != expectedHash {
			t.Fatalf("Expected page hash to be %s, got %s", expectedHash, pageCrawled.Hash)
		}

		if err != nil {
			t.Fatalf("Error creating repository: %v", err)
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
			t.Errorf("Expected all mocks to be called")
		}
	})

	t.Run("handles 404", func(t *testing.T) {
		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo, nil)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		service := service.NewService(htmlService, pageService, nil, nil, testutil.NewTestLogger())

		defer gock.Off()

		gock.New("http://example.com").
			Get("/").
			Reply(404)

		_, err := service.Crawl("test1", "http://example.com")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("only crawls html content type", func(t *testing.T) {
		defer gock.Off()

		htmlRepo := htmlRepo.NewRepository()
		htmlBackupService := htmlBackupService.NewService(html.NewClient("http://storage:8080"))
		htmlService := htmlService.NewService(htmlRepo, htmlBackupService)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		indexService := indexService.NewService(pageService, htmlService, peerService, keywordService, testutil.NewTestLogger())

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(htmlService, pageService, indexService, api, testutil.NewTestLogger())

		expectedData := `<html><head><title>Example</title></head><body><h1>Hello, World! <a href="/about">About us</a></h1><p>Welcome to my website.</p></body></html>`

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			SetHeader("Content-Type", "application/atom+xml")

		pageCrawled, err := service.Crawl("test1", "http://example.com")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if pageCrawled != nil {
			t.Fatalf("Expected page to be nil")
		}

		if !gock.IsDone() {
			t.Errorf("Expected all mocks to be called")
		}
	})
}

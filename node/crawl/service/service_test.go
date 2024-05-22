package service_test

import (
	"crawlquery/node/crawl/service"
	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	peerService "crawlquery/node/peer/service"

	indexService "crawlquery/node/index/service"

	"crawlquery/pkg/client/api"
	"crawlquery/pkg/testutil"
	"testing"

	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	t.Run("can crawl a page", func(t *testing.T) {

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)
		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(pageService, htmlService, keywordService, peerService, testutil.NewTestLogger())

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(htmlService, pageService, indexService, api, testutil.NewTestLogger())

		defer gock.Off()

		expectedData := "<html><head><title>Example</title></head><body><h1>Hello, World!</h1></body></html>"

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData)

		err := service.Crawl("test1", "http://example.com")

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
		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)
		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(pageService, htmlService, keywordService, peerService, testutil.NewTestLogger())

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(htmlService, pageService, indexService, api, testutil.NewTestLogger())

		defer gock.Off()

		expectedData := `<html><head><title>Example</title></head><body><h1>Hello, World! <a href="http://example.com/about">About us</a></h1></body></html>`

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			SetHeader("Content-Type", "text/html")

		gock.New("http://localhost:8080").
			Post("/crawl-jobs").
			JSON(`{"url":"http://example.com/about"}`).
			Reply(200)

		err := service.Crawl("test1", "http://example.com")

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

	t.Run("handles 404", func(t *testing.T) {
		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		service := service.NewService(htmlService, pageService, nil, nil, testutil.NewTestLogger())

		defer gock.Off()

		gock.New("http://example.com").
			Get("/").
			Reply(404)

		err := service.Crawl("test1", "http://example.com")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

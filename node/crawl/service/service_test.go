package service_test

import (
	"strings"
	"testing"

	crawlService "crawlquery/node/crawl/service"
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

	"github.com/gocolly/colly/v2"
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

	crawlService := crawlService.NewService(htmlSvc, pageSvc, indexSvc, apiClient, testutil.NewTestLogger())

	return crawlService, htmlRepository, pageRepository
}

func TestCrawl(t *testing.T) {
	t.Run("can crawl a page", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://storage:8080").Post("/pages").Reply(201)

		gock.New("http://example.com:9292/robots.txt").
			Reply(200).
			// allow all
			BodyString("User-agent: *\nAllow: /")

		service, htmlRepo, pageRepo := setupServices()

		expectedData := "<html><head><title>Example</title></head><body><h1>Hello, World!</h1><p>Welcome to my example website.</p></body></html>"

		gock.New("http://example.com:9292").
			Get("/").
			Reply(200).
			BodyString(expectedData).Header.Set("Content-Type", "text/html")

		_, _, err := service.Crawl("test1", "http://example.com:9292")

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
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

		if page.URL != "http://example.com:9292" {
			t.Fatalf("Expected page URL to be 'http://example.com:9292', got '%s'", page.URL)
		}

		if !gock.IsDone() {
			t.Errorf("Expected all mocks to be called")
		}
	})

	t.Run("creates crawl jobs for links", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://storage:8080").Post("/pages").Reply(201)

		gock.New("http://example.com/robots.txt").
			Reply(200).
			// allow all
			BodyString("User-agent: *\nAllow: /")

		service, htmlRepo, pageRepo := setupServices()

		expectedData := `<html><head><title>Example</title></head><body><h1>Hello, World! <a href="http://example.com/about">About us</a></h1><p>Welcome to my website.</p></body></html>`

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			SetHeader("Content-Type", "text/html")

		hash, links, err := service.Crawl("test1", "http://example.com")

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
		}

		expectedHash := util.Sha256Hex32([]byte(expectedData))

		if hash != expectedHash {
			t.Fatalf("Expected page hash to be %s, got %s", expectedHash, hash)
		}

		if len(links) != 1 {
			t.Fatalf("Expected 1 link, got %d", len(links))
		}

		if links[0] != "http://example.com/about" {
			t.Fatalf("Expected link to be 'http://example.com/about', got '%s'", links[0])
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

		gock.New("http://storage:8080").Post("/pages").Reply(201)

		service, htmlRepo, pageRepo := setupServices()

		expectedData := `<html><head><title>Example</title></head><body><h1>Hello, World! <a href="/about">About us</a></h1><p>Welcome to my website.</p></body></html>`

		gock.New("http://example.com/robots.txt").
			Reply(200).
			// allow all
			BodyString("User-agent: *\nAllow: /")

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			SetHeader("Content-Type", "text/html")

		hash, links, err := service.Crawl("test1", "http://example.com")

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
		}

		expectedHash := util.Sha256Hex32([]byte(expectedData))

		if hash != expectedHash {
			t.Fatalf("Expected page hash to be %s, got %s", expectedHash, hash)
		}

		if len(links) != 1 {
			t.Fatalf("Expected 1 link, got %d", len(links))
		}

		if links[0] != "http://example.com/about" {
			t.Fatalf("Expected link to be 'http://example.com/about', got '%s'", links[0])
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
		defer gock.Off()

		service, _, _ := setupServices()

		gock.New("http://example.com").
			Get("/").
			Reply(404)

		_, _, err := service.Crawl("test1", "http://example.com")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("only crawls html content type", func(t *testing.T) {
		defer gock.Off()

		service, _, _ := setupServices()

		expectedData := `<html><head><title>Example</title></head><body><h1>Hello, World! <a href="/about">About us</a></h1><p>Welcome to my website.</p></body></html>`

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData).
			SetHeader("Content-Type", "application/atom+xml")

		_, _, err := service.Crawl("test1", "http://example.com")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if !gock.IsDone() {
			t.Errorf("Expected all mocks to be called")
		}
	})

	t.Run("simulate robots.txt failure", func(t *testing.T) {

		service, _, _ := setupServices()

		gock.New("http://example.com").
			Get("/robots.txt").
			Reply(200).
			BodyString("User-agent: *\nDisallow: /")

		gock.New("http://example.com").
			Get("/").
			Reply(400).
			BodyString("<html><head><title>Example</title></head><body><h1>Hello, World!</h1><p>Welcome to my example website.</p></body></html>").
			SetHeader("Content-Type", "text/html")

		_, _, err := service.Crawl("test1", "http://example.com")

		if err != colly.ErrRobotsTxtBlocked {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("does not follow redirects", func(t *testing.T) {
		defer gock.Off()

		service, _, _ := setupServices()

		gock.New("http://exampleredirect.com").
			Get("/").
			Reply(301).
			SetHeader("Location", "http://example.com/redirect")

		_, _, err := service.Crawl("test1", "http://exampleredirect.com")

		if !strings.Contains(err.Error(), "301") {
			t.Errorf("Expected error to contain '301', got '%s'", err.Error())
		}
	})
}

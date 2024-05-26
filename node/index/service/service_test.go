package service_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/index/service"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	peerService "crawlquery/node/peer/service"

	"testing"
)

func TestIndex(t *testing.T) {

	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo, nil)

	pageRepo.Save("page1", &domain.Page{
		ID:  "page1",
		URL: "http://example.com",
	})

	htmlRepo := htmlRepo.NewRepository()
	htmlService := htmlService.NewService(htmlRepo, nil)

	keywordRepo := keywordRepo.NewRepository()
	keywordService := keywordService.NewService(keywordRepo)

	htmlRepo.Save("page1", []byte(`
		<html>
			<head>
				<title>Test Page</title>
				<meta name="og:description" content="This is a test page">
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page, with some good english words.</p>
			</body>
		</html>
	`))

	peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

	logger := testutil.NewTestLogger()

	s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

	err := s.Index("page1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	page, err := pageRepo.Get("page1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if page == nil {
		t.Fatalf("Expected page to be found, got nil")
	}

	if page.ID != "page1" {
		t.Fatalf("Expected page ID to be page1, got %s", page.ID)
	}

	if page.URL != "http://example.com" {
		t.Fatalf("Expected URL to be http://example.com, got %s", page.URL)
	}

	if page.Title != "Test Page" {
		t.Fatalf("Expected title to be Test Page, got %s", page.Title)
	}

	if page.Description != "This is a test page, with some good english words." {
		t.Fatalf("Expected meta description to be This is a test page, got %s", page.Description)
	}

	if page.LastIndexedAt.IsZero() {
		t.Fatalf("Expected last indexed at to be set, got zero")
	}
}

func TestGetIndex(t *testing.T) {
	t.Run("can get index", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "http://example.com",
		})

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo, nil)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		html := []byte(`
		<html>
			<head>
				<title>Test Page</title>
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page, with some good english words.</p>
			</body>
		</html>
	`)
		htmlRepo.Save("page1", html)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		logger := testutil.NewTestLogger()

		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		err := s.Index("page1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		page, err := s.GetIndex("page1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if page == nil {
			t.Fatalf("Expected page to be found, got nil")
		}

		if page.ID != "page1" {
			t.Fatalf("Expected page ID to be page1, got %s", page.ID)
		}

		if page.URL != "http://example.com" {
			t.Fatalf("Expected URL to be http://example.com, got %s", page.URL)
		}

		if page.Title != "Test Page" {
			t.Fatalf("Expected title to be Test Page, got %s", page.Title)
		}

		expectedHash := util.Sha256Hex32(html)

		if page.Hash != expectedHash {
			t.Fatalf("Expected hash to be %s, got %s", expectedHash, page.Hash)
		}
	})
}

func TestSearch(t *testing.T) {
	t.Run("can search for keyword", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "http://example.com",
		})

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo, nil)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		html := []byte(`
		<html>
			<head>
				<title>Test Page</title>
				<meta name="description" content="This is a test page">
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page, with some good english words.</p>
			</body>
		</html>
	`)
		htmlRepo.Save("page1", html)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		logger := testutil.NewTestLogger()

		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		err := s.Index("page1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		results, err := s.Search("test page")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].PageID != "page1" {
			t.Fatalf("Expected page ID to be page1, got %s", results[0].PageID)
		}

		if results[0].Score < 100 {
			t.Fatalf("Expected score to be more, got %f", results[0].Score)
		}

		if results[0].Page.URL != "http://example.com" {
			t.Fatalf("Expected URL to be http://example.com, got %s", results[0].Page.URL)
		}

		if results[0].Page.Title != "Test Page" {
			t.Fatalf("Expected title to be Test Page, got %s", results[0].Page.Title)
		}

		if results[0].Page.Description != "This is a test page" {
			t.Fatalf("Expected meta description to be This is a test page, got %s", results[0].Page.Description)
		}

		expectedHash := util.Sha256Hex32(html)

		if results[0].Page.Hash != expectedHash {
			t.Fatalf("Expected hash to be %s, got %s", expectedHash, results[0].Page.Hash)
		}
	})

	t.Run("applies domain signal", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "https://youtube.com",
		})

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo, nil)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		htmlRepo.Save("page1", []byte(`
		<html>
			<head>
				<title>Test Page</title>
				<meta name="description" content="This is a test page">
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page, with some good english words.</p>
			</body>
		</html>
	`))

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		logger := testutil.NewTestLogger()

		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		err := s.Index("page1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		results, err := s.Search("youtube.com")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].Score < 1000 {
			t.Fatalf("Expected score to be 1000 or more, got %f", results[0].Score)
		}
	})

	t.Run("sets signal breakdown", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "http://example.com",
		})

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo, nil)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		html := []byte(`
		<html>
			<head>
				<title>Test Page</title>
				<meta name="description" content="This is a test page">
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page, with some good english words.</p>
			</body>
		</html>
	`)
		htmlRepo.Save("page1", html)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		logger := testutil.NewTestLogger()

		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		err := s.Index("page1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		results, err := s.Search("example")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if len(results[0].Signals) != 3 {
			t.Fatalf("Expected 3 signals, got %d", len(results[0].Signals))
		}

		if results[0].Signals["domain"]["domain"] != 2000 {
			t.Fatalf("Expected domain signal to be 2000, got %f", results[0].Signals["domain"]["domain"])
		}
	})
}

func TestApplyPageUpdatedEvent(t *testing.T) {
	t.Run("can apply index event", func(t *testing.T) {

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo, nil)

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		service := service.NewService(pageService, htmlService, peerService, keywordService, testutil.NewTestLogger())

		page := &domain.Page{
			URL:         "http://example.com",
			ID:          "page1",
			Title:       "Example",
			Description: "An example page",
		}

		event := &domain.PageUpdatedEvent{
			Page: page,
		}

		err := service.ApplyPageUpdatedEvent(event)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		page, err = pageService.Get("page1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if page == nil {
			t.Fatalf("Expected page to be found, got nil")
		}

		if page.ID != "page1" {
			t.Fatalf("Expected page ID to be page1, got %s", page.ID)
		}

		if page.Title != "Example" {
			t.Fatalf("Expected title to be Example, got %s", page.Title)
		}

		if page.Description != "An example page" {
			t.Fatalf("Expected meta description to be An example page, got %s", page.Description)
		}

	})

}

func TestHash(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo, nil)

	pageRepo.Save("page1", &domain.Page{
		ID:  "page1",
		URL: "http://example.com",
	})

	htmlRepo := htmlRepo.NewRepository()
	htmlService := htmlService.NewService(htmlRepo, nil)

	keywordRepo := keywordRepo.NewRepository()
	keywordService := keywordService.NewService(keywordRepo)

	htmlRepo.Save("page1", []byte(`
		<html>
			<head>
				<title>Test Page</title>
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page, with some good english words.</p>
			</body>
		</html>
	`))

	peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

	logger := testutil.NewTestLogger()

	s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

	err := s.Index("page1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	pageHash, err := s.Hash()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	hash, err := pageRepo.GetHash("page1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pageHash != util.Sha256Hex32([]byte(hash)) {
		t.Fatalf("Expected hash to be %s, got %s", util.Sha256Hex32([]byte(hash)), pageHash)
	}
}

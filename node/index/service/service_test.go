package service_test

import (
	"crawlquery/node/index/service"
	"crawlquery/pkg/domain"
	"crawlquery/pkg/testutil"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	"testing"
)

func TestMakePostings(t *testing.T) {
	s := service.NewService(nil, nil, nil, nil)

	postings := s.MakePostings(&domain.Page{
		ID: "page1",
	}, []string{"test1", "test2", "test2"})

	if len(postings) != 2 {
		t.Fatalf("Expected 2 postings, got %d", len(postings))
	}

	if postings["test1"].Frequency != 1 {
		t.Fatalf("Expected frequency to be 1, got %d", postings["test1"].Frequency)
	}

	if postings["test2"].Frequency != 2 {
		t.Fatalf("Expected frequency to be 2, got %d", postings["test2"].Frequency)
	}

	if postings["test1"].Positions[0] != 0 {
		t.Fatalf("Expected position to be 0, got %d", postings["test1"].Positions[0])
	}

	if postings["test2"].Positions[0] != 1 {
		t.Fatalf("Expected position to be 1, got %d", postings["test2"].Positions[0])
	}

	if postings["test2"].Positions[1] != 2 {
		t.Fatalf("Expected position to be 2, got %d", postings["test2"].Positions[1])
	}
}

func TestIndex(t *testing.T) {

	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo)

	pageRepo.Save("page1", &domain.Page{
		ID:  "page1",
		URL: "http://example.com",
	})

	htmlRepo := htmlRepo.NewRepository()
	htmlService := htmlService.NewService(htmlRepo)

	htmlRepo.Save("page1", []byte(`
		<html>
			<head>
				<title>Test Page</title>
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page</p>
			</body>
		</html>
	`))

	keywordRepo := keywordRepo.NewRepository()
	keywordService := keywordService.NewService(keywordRepo)

	logger := testutil.NewTestLogger()

	s := service.NewService(pageService, htmlService, keywordService, logger)

	err := s.Index("page1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	postings, err := keywordRepo.GetPostings("test")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(postings) != 1 {
		t.Fatalf("Expected 1 posting, got %d", len(postings))
	}

	if postings[0].PageID != "page1" {
		t.Fatalf("Expected page ID to be page1, got %s", postings[0].PageID)
	}

	if postings[0].Frequency != 2 {
		t.Fatalf("Expected frequency to be 2, got %d", postings[0].Frequency)
	}
}

func TestSearch(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo)

	pageRepo.Save("page1", &domain.Page{
		ID:  "page1",
		URL: "http://example.com",
	})

	htmlRepo := htmlRepo.NewRepository()
	htmlService := htmlService.NewService(htmlRepo)

	htmlRepo.Save("page1", []byte(`
		<html>
			<head>
				<title>Test Page</title>
				<meta name="description" content="This is a test page">
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This is a test page</p>
			</body>
		</html>
	`))

	keywordRepo := keywordRepo.NewRepository()
	keywordService := keywordService.NewService(keywordRepo)

	logger := testutil.NewTestLogger()

	s := service.NewService(pageService, htmlService, keywordService, logger)

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

	if results[0].Score != 4 {
		t.Fatalf("Expected score to be 4, got %f", results[0].Score)
	}

	if results[0].Page.URL != "http://example.com" {
		t.Fatalf("Expected URL to be http://example.com, got %s", results[0].Page.URL)
	}

	if results[0].Page.Title != "Test Page" {
		t.Fatalf("Expected title to be Test Page, got %s", results[0].Page.Title)
	}

	if results[0].Page.MetaDescription != "This is a test page" {
		t.Fatalf("Expected meta description to be This is a test page, got %s", results[0].Page.MetaDescription)
	}
}

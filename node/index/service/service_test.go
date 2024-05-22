package service_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/index/service"
	sharedDomain "crawlquery/pkg/domain"
	"crawlquery/pkg/testutil"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	peerService "crawlquery/node/peer/service"

	"testing"
)

func TestMakePostings(t *testing.T) {
	s := service.NewService(nil, nil, nil, nil, nil)

	postings := s.MakePostings(&sharedDomain.Page{
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

	pageRepo.Save("page1", &sharedDomain.Page{
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

	peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

	logger := testutil.NewTestLogger()

	s := service.NewService(pageService, htmlService, keywordService, peerService, logger)

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

	pageRepo.Save("page1", &sharedDomain.Page{
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

	peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

	logger := testutil.NewTestLogger()

	s := service.NewService(pageService, htmlService, keywordService, peerService, logger)

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

func TestApplyIndexEvent(t *testing.T) {
	t.Run("can apply index event", func(t *testing.T) {

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

		service := service.NewService(pageService, htmlService, keywordService, peerService, testutil.NewTestLogger())

		page := &sharedDomain.Page{
			URL:             "http://example.com",
			ID:              "page1",
			Title:           "Example",
			MetaDescription: "An example page",
		}

		event := &domain.IndexEvent{
			Page: page,
			Keywords: map[string]*domain.Posting{
				"keyword1": {
					PageID:    "page1",
					Frequency: 1,
					Positions: []int{1},
				},
			},
		}

		err := service.ApplyIndexEvent(event)

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

		if page.MetaDescription != "An example page" {
			t.Fatalf("Expected meta description to be An example page, got %s", page.MetaDescription)
		}

		postings, err := keywordService.GetPostings("keyword1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(postings) != 1 {
			t.Fatalf("Expected 1 posting, got %d", len(postings))
		}

		if postings[0].PageID != "page1" {
			t.Fatalf("Expected page id to be page1, got %s", postings[0].PageID)
		}

		if postings[0].Frequency != 1 {
			t.Fatalf("Expected frequency to be 1, got %d", postings[0].Frequency)
		}
	})

	t.Run("can remove old postings", func(t *testing.T) {
		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

		service := service.NewService(pageService, htmlService, keywordService, peerService, testutil.NewTestLogger())

		page := &sharedDomain.Page{
			URL:             "http://example.com",
			ID:              "page1",
			Title:           "Example",
			MetaDescription: "An example page",
		}

		event := &domain.IndexEvent{
			Page: page,
			Keywords: map[string]*domain.Posting{
				"keyword1": {
					PageID:    "page1",
					Frequency: 1,
					Positions: []int{1},
				},
			},
		}

		err := service.ApplyIndexEvent(event)

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

		postings, err := keywordService.GetPostings("keyword1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(postings) != 1 {
			t.Fatalf("Expected 1 posting, got %d", len(postings))
		}

		if postings[0].PageID != "page1" {
			t.Fatalf("Expected page id to be page1, got %s", postings[0].PageID)
		}

		if postings[0].Frequency != 1 {
			t.Fatalf("Expected frequency to be 1, got %d", postings[0].Frequency)
		}

		event = &domain.IndexEvent{
			Page: page,
			Keywords: map[string]*domain.Posting{
				"keyword1": {
					PageID:    "page1",
					Frequency: 2,
					Positions: []int{1},
				},
			},
		}

		err = service.ApplyIndexEvent(event)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		postings, err = keywordService.GetPostings("keyword1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(postings) != 1 {
			t.Fatalf("Expected 1 posting, got %d", len(postings))
		}

		if postings[0].PageID != "page1" {
			t.Fatalf("Expected page id to be page1, got %s", postings[0].PageID)
		}

		if postings[0].Frequency != 2 {
			t.Fatalf("Expected frequency to be 2, got %d", postings[0].Frequency)
		}
	})
}

func TestHash(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(pageRepo)

	pageRepo.Save("page1", &sharedDomain.Page{
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

	peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

	logger := testutil.NewTestLogger()

	s := service.NewService(pageService, htmlService, keywordService, peerService, logger)

	err := s.Index("page1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	pageHash, keywordHash, combinedHash, err := s.Hash()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pageHash != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatalf("Expected page hash to be e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855, got %s", pageHash)
	}

	if keywordHash != "ee39e266c3ed88a122d41ec182d5aaf206d20dd5db91b27ffba4d0893c8e16fc" {
		t.Fatalf("Expected keyword hash to be ee39e266c3ed88a122d41ec182d5aaf206d20dd5db91b27ffba4d0893c8e16fc, got %s", keywordHash)
	}

	if combinedHash != "922ffc6f07eb2bab44239b8bebce271cb982c3720ff6cd31d679e2c2c20768b5" {
		t.Fatalf("Expected combined hash to be 922ffc6f07eb2bab44239b8bebce271cb982c3720ff6cd31d679e2c2c20768b5, got %s", combinedHash)
	}
}

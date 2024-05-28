package service_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/index/service"
	"crawlquery/pkg/testutil"
	"fmt"
	"time"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	peerService "crawlquery/node/peer/service"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	"testing"

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

	htmlRepo := htmlRepo.NewRepository()
	htmlService := htmlService.NewService(htmlRepo, nil)

	peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

	keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
	keywordService := keywordService.NewService(keywordOccurrenceRepo)

	return pageRepo, pageService, htmlRepo, htmlService, peerService, keywordService
}

func TestIndex(t *testing.T) {

	t.Run("indexes a page", func(t *testing.T) {
		pageRepo, pageService, htmlRepo, htmlService, peerService, keywordService := setupTestRepos()

		logger := testutil.NewTestLogger()
		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "http://example.com",
		})

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

		checkOccurrences := map[domain.Keyword]domain.KeywordOccurrence{
			"page": {
				PageID:    "page1",
				Frequency: 2,
				Positions: []int{3, 4},
			},
		}

		for keyword, check := range checkOccurrences {
			matches, err := keywordService.GetKeywordMatches([]domain.Keyword{keyword})
			if err != nil {
				t.Fatalf("Error getting occurrences: %v", err)
			}

			if len(matches) != 1 {
				t.Fatalf("Expected 1 match, got %d", len(matches))
			}

			if matches[0].Occurrences[0].PageID != check.PageID {
				t.Fatalf("Expected page ID %s, got %s", check.PageID, matches[0].Occurrences[0].PageID)
			}
		}

	})

	t.Run("trunacates keywords when more than 1500 found", func(t *testing.T) {
		pageRepo, pageService, htmlRepo, htmlService, peerService, keywordService := setupTestRepos()

		logger := testutil.NewTestLogger()
		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "http://example.com",
		})

		start := []byte(`
		<html>
			<head>
				<title>Test Page</title>
				<meta name="og:description" content="This is a test page">
			</head>

			<body>
				<h1>Test Page</h1>
				<p>This`)

		for i := 0; i < 1500; i++ {
			start = append(start, []byte(fmt.Sprintf(" Test%d ", i))...)
		}

		start = append(start, []byte(`</p>
			</body>
		</html>
	`)...)

		htmlRepo.Save("page1", start)

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

		if page.LastIndexedAt.IsZero() {
			t.Fatalf("Expected last indexed at to be set, got zero")
		}

		count, err := keywordService.Count()

		if err != nil {
			t.Fatalf("Error counting keywords: %v", err)
		}

		if count != 1500 {
			t.Fatalf("Expected 1500 keywords, got %d", count)
		}
	})

	t.Run("sends page updated event", func(t *testing.T) {
		pageRepo, pageService, htmlRepo, htmlService, peerService, keywordService := setupTestRepos()

		logger := testutil.NewTestLogger()
		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "http://example.com",
		})

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

		peerService.AddPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "node1.cluster.com",
			Port:     8080,
		})

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Post("/event").
			Reply(200)

		err := s.Index("page1")

		time.Sleep(50 * time.Millisecond)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

func TestApplyPageUpdatedEvent(t *testing.T) {
	t.Run("updates a page", func(t *testing.T) {
		pageRepo, pageService, htmlRepo, htmlService, peerService, keywordService := setupTestRepos()

		logger := testutil.NewTestLogger()
		s := service.NewService(pageService, htmlService, peerService, keywordService, logger)

		pageRepo.Save("page1", &domain.Page{
			ID:  "page1",
			URL: "http://example.com",
		})

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

		page.Title = "New Title"
		page.Description = "New Description"

		err = s.ApplyPageUpdatedEvent(&domain.PageUpdatedEvent{
			Page: page,
			KeywordOccurrences: map[domain.Keyword]domain.KeywordOccurrence{
				"keyword1": {
					PageID:    "page1",
					Frequency: 1,
					Positions: []int{1, 2, 3},
				},
			},
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		page, err = pageRepo.Get("page1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if page.Title != "New Title" {
			t.Fatalf("Expected title to be New Title, got %s", page.Title)
		}

		if page.Description != "New Description" {
			t.Fatalf("Expected description to be New Description, got %s", page.Description)
		}

		checkOccurrences, err := keywordService.GetForPageID("page1")

		if err != nil {
			t.Fatalf("Error getting occurrences: %v", err)
		}

		if len(checkOccurrences) != 1 {
			t.Fatalf("Expected 1 occurrence, got %d", len(checkOccurrences))
		}
	})
}

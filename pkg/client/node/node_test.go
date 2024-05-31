package node_test

import (
	"crawlquery/node/dto"
	"crawlquery/pkg/client/node"
	"time"

	"testing"

	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		defer gock.Off()

		expectedRes := &dto.CrawlResponse{
			ContentHash: "hash",
			Links: []string{
				"http://example.com",
			},
		}

		gock.New("http://node.com").
			Post("/crawl").
			Reply(200).
			JSON(expectedRes)

		node := node.NewClient(
			node.WithHostname("node.com"),
			node.WithPort(80),
		)

		res, err := node.Crawl("page1", "http://example.com")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if res.ContentHash != expectedRes.ContentHash {
			t.Fatalf("Expected %s, got %s", expectedRes.ContentHash, res)
		}

		if res.Links[0] != expectedRes.Links[0] {
			t.Fatalf("Expected %s, got %s", expectedRes.Links[0], res.Links[0])
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

func TestIndex(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		defer gock.Off()

		expectedRes := &dto.IndexResponse{
			Success: true,
		}

		gock.New("http://node.com").
			Post("/index").
			JSON(&dto.IndexRequest{
				PageID:      "page1",
				URL:         "http://example.com",
				ContentHash: "hash",
			}).
			Reply(200).
			JSON(expectedRes)

		node := node.NewClient(
			node.WithHostname("node.com"),
			node.WithPort(80),
		)

		err := node.Index("page1", "http://example.com", "hash")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

func TestGetIndexMetas(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		defer gock.Off()

		expectedRes := &dto.GetIndexMetasResponse{
			IndexMetas: []dto.IndexMeta{
				{
					PageID:        "page1",
					LastIndexedAt: time.Now(),
				},
			},
		}

		gock.New("http://node.com").
			Post("/repair/get-index-metas").
			JSON(&dto.GetIndexMetasRequest{
				PageIDs: []string{"page1"},
			}).
			Reply(200).
			JSON(expectedRes)

		node := node.NewClient(
			node.WithHostname("node.com"),
			node.WithPort(80),
		)

		indexMetas, err := node.GetIndexMetas([]string{"page1"})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if indexMetas[0].PageID != expectedRes.IndexMetas[0].PageID {
			t.Fatalf("Expected %s, got %s", expectedRes.IndexMetas[0].PageID, indexMetas[0].PageID)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

func TestGetAllIndexMetas(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		defer gock.Off()

		expectedRes := &dto.GetIndexMetasResponse{
			IndexMetas: []dto.IndexMeta{
				{
					PageID:        "page1",
					LastIndexedAt: time.Now(),
				},
			},
		}

		gock.New("http://node.com").
			Get("/repair/get-all-index-metas").
			Reply(200).
			JSON(expectedRes)

		node := node.NewClient(
			node.WithHostname("node.com"),
			node.WithPort(80),
		)

		indexMetas, err := node.GetAllIndexMetas()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if indexMetas[0].PageID != expectedRes.IndexMetas[0].PageID {
			t.Fatalf("Expected %s, got %s", expectedRes.IndexMetas[0].PageID, indexMetas[0].PageID)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

func TestGetPageDumps(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		defer gock.Off()

		expectedRes := &dto.GetPageDumpsResponse{
			PageDumps: []dto.PageDump{
				{
					PageID: "page1",
					Page: dto.Page{
						ID:          "page1",
						URL:         "http://example.com",
						Title:       "Example",
						Description: "Description",
					},
					KeywordOccurrences: map[string]dto.KeywordOccurrence{
						"keyword1": {
							PageID:    "page1",
							Frequency: 1,
							Positions: []int{1},
						},
					},
				},
			},
		}

		gock.New("http://node.com").
			Post("/repair/get-page-dumps").
			JSON(&dto.GetPageDumpsRequest{
				PageIDs: []string{"page1"},
			}).
			Reply(200).
			JSON(expectedRes)

		node := node.NewClient(
			node.WithHostname("node.com"),
			node.WithPort(80),
		)

		res, err := node.GetPageDumps([]string{"page1"})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if res[0].Page.ID != expectedRes.PageDumps[0].Page.ID {
			t.Fatalf("Expected %s, got %s", expectedRes.PageDumps[0].Page.ID, res[0].Page.ID)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

package node_test

import (
	"crawlquery/node/dto"
	"crawlquery/pkg/client/api/node"

	"testing"

	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		defer gock.Off()

		expectedRes := &dto.CrawlResponse{
			Page: &dto.Page{
				ID:   "page1",
				URL:  "http://example.com",
				Hash: "hash",
			},
		}

		gock.New("http://node.com").
			Post("/crawl").
			Reply(200).
			JSON(expectedRes)

		node := node.NewClient("http://node.com")

		res, err := node.Crawl("page1", "http://example.com")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if res.Hash != expectedRes.Page.Hash {
			t.Fatalf("Expected %s, got %s", expectedRes.Page.Hash, res)
		}

		if !gock.IsDone() {
			t.Fatalf("Expected all mocks to be called")
		}
	})
}

package service_test

import (
	"crawlquery/node/domain"
	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	"testing"

	dumpService "crawlquery/node/dump/service"
)

func TestKeyword(t *testing.T) {
	t.Run("can get keyword dump", func(t *testing.T) {
		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		dumpService := dumpService.NewService(nil, keywordService)

		data, err := dumpService.Keyword()

		if err != nil {
			t.Fatalf("Error getting keyword dump: %v", err)
		}

		if string(data) != "{}" {
			t.Fatalf("Expected empty object, got %s", string(data))
		}

		keywordService.SavePostings(map[string]*domain.Posting{
			"hello": {PageID: "1", Frequency: 1, Positions: []int{1}},
		})

		data, err = dumpService.Keyword()

		if err != nil {
			t.Fatalf("Error getting keyword dump: %v", err)
		}

		if string(data) != `{"hello":[{"page_id":"1","frequency":1,"positions":[1]}]}` {
			t.Fatalf("Expected keyword dump to be 'hello', got %s", string(data))
		}
	})
}

func TestPage(t *testing.T) {
	t.Run("can get page dump", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		dumpService := dumpService.NewService(pageService, nil)

		data, err := dumpService.Page()

		if err != nil {
			t.Fatalf("Error getting page dump: %v", err)
		}

		if string(data) != "{}" {
			t.Fatalf("Expected empty object, got %s", string(data))
		}

		pageService.Create("1", "http://example.com")

		data, err = dumpService.Page()

		if err != nil {
			t.Fatalf("Error getting page dump: %v", err)
		}

		if string(data) != `{"1":{"id":"1","url":"http://example.com","title":"","meta_description":""}}` {
			t.Fatalf("Expected page dump to be '1', got %s", string(data))
		}
	})
}

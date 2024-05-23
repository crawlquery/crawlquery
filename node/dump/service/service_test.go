package service_test

import (
	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	"testing"

	dumpService "crawlquery/node/dump/service"
)

func TestPage(t *testing.T) {
	t.Run("can get page dump", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		dumpService := dumpService.NewService(pageService)

		data, err := dumpService.Page()

		if err != nil {
			t.Fatalf("Error getting page dump: %v", err)
		}

		if string(data) != "{}" {
			t.Fatalf("Expected empty object, got %s", string(data))
		}

		pageRepo.Save("1", &domain.Page{
			ID:       "1",
			URL:      "http://example.com",
			Hash:     "1",
			Keywords: []string{"test", "page"},
		})

		data, err = dumpService.Page()

		if err != nil {
			t.Fatalf("Error getting page dump: %v", err)
		}

		if string(data) != `{"1":{"id":"1","hash":"1","url":"http://example.com","title":"","meta_description":"","keywords":["test","page"]}}` {
			t.Fatalf("Expected page dump to be '1', got %s", string(data))
		}
	})
}

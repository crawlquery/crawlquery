package service_test

import (
	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"
	"encoding/json"

	"testing"

	dumpService "crawlquery/node/dump/service"
)

func TestPage(t *testing.T) {
	t.Run("can get page dump", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		dumpService := dumpService.NewService(pageService)

		data, err := dumpService.Page()

		if err != nil {
			t.Fatalf("Error getting page dump: %v", err)
		}

		if string(data) != "{}" {
			t.Fatalf("Expected empty object, got %s", string(data))
		}

		pageRepo.Save("1", &domain.Page{
			ID:   "1",
			URL:  "http://example.com",
			Hash: "1",
		})

		data, err = dumpService.Page()

		if err != nil {
			t.Fatalf("Error getting page dump: %v", err)
		}

		var slicePages map[string]*domain.Page

		err = json.Unmarshal(data, &slicePages)

		if err != nil {
			t.Fatalf("Error unmarshalling page dump: %v", err)
		}

		if len(slicePages) != 1 {
			t.Fatalf("Expected 1 page, got %d", len(slicePages))
		}

		if slicePages["1"].ID != "1" {
			t.Fatalf("Expected page ID to be '1', got '%s'", slicePages["1"].ID)
		}

		if slicePages["1"].URL != "http://example.com" {
			t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", slicePages["1"].URL)
		}

		if slicePages["1"].Hash != "1" {
			t.Fatalf("Expected page Hash to be '1', got '%s'", slicePages["1"].Hash)
		}

	})
}

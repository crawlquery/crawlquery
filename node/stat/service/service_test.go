package service_test

import (
	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"
	"encoding/json"
	"testing"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	dumpService "crawlquery/node/dump/service"

	statService "crawlquery/node/stat/service"
)

func TestInfo(t *testing.T) {
	t.Run("returns stat info", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		dumpService := dumpService.NewService(pageService)

		keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
		keywordService := keywordService.NewService(keywordOccurrenceRepo)

		keywordOccurrenceRepo.Add(domain.Keyword("example"), domain.KeywordOccurrence{
			PageID:    "1",
			Frequency: 1,
			Positions: []int{1},
		})

		statService := statService.NewService(pageService, keywordService, dumpService)

		pages := map[string]*domain.Page{
			"1": {
				ID:          "1",
				URL:         "http://example.com",
				Title:       "Example Domain",
				Description: "",
			},
			"2": {
				ID:          "2",
				URL:         "http://example.com",
				Title:       "Example Domain",
				Description: "",
			},
			"3": {
				ID:          "3",
				URL:         "http://example.com",
				Title:       "Example Domain",
				Description: "",
			},
		}

		encoded, err := json.Marshal(pages)

		if err != nil {
			t.Fatalf("error marshalling pages: %v", err)
		}

		for _, p := range pages {
			err = pageRepo.Save(p.ID, p)
			if err != nil {
				t.Fatalf("error saving page: %v", err)
			}
		}

		info, err := statService.Info()

		if err != nil {
			t.Fatalf("error getting stat info: %v", err)
		}

		if info.TotalPages != 3 {
			t.Errorf("expected 3 pages, got %d", info.TotalPages)
		}

		if info.TotalKeywords != 1 {
			t.Errorf("expected 1 keyword, got %d", info.TotalKeywords)
		}

		if info.SizeOfPages != len(encoded) {
			t.Errorf("expected %d page size, got %d", len(encoded), info.SizeOfPages)
		}
	})
}

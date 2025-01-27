package handler_test

import (
	"bytes"
	"crawlquery/node/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	dumpService "crawlquery/node/dump/service"

	statService "crawlquery/node/stat/service"

	statHandler "crawlquery/node/stat/handler"

	"github.com/gin-gonic/gin"
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

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/info", bytes.NewReader(encoded))

		statHandler := statHandler.NewHandler(statService)

		statHandler.Info(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		var info domain.StatInfo

		err = json.NewDecoder(w.Body).Decode(&info)

		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		if info.TotalPages != 3 {
			t.Errorf("expected 3 pages, got %d", info.TotalPages)
		}

		if info.SizeOfPages != len(encoded) {
			t.Errorf("expected %d index size, got %d", len(encoded), info.SizeOfPages)
		}
	})
}

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

	dumpService "crawlquery/node/dump/service"

	statService "crawlquery/node/stat/service"

	statHandler "crawlquery/node/stat/handler"

	"github.com/gin-gonic/gin"
)

func TestInfo(t *testing.T) {
	t.Run("returns stat info", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		dumpService := dumpService.NewService(pageService)

		statService := statService.NewService(pageService, dumpService)

		pages := map[string]*domain.Page{
			"1": {
				ID:          "1",
				URL:         "http://example.com",
				Title:       "Example Domain",
				Description: "",
				Phrases:     [][]string{{"example", "domain"}},
			},
			"2": {
				ID:          "2",
				URL:         "http://example.com",
				Title:       "Example Domain",
				Description: "",
				Phrases:     [][]string{{"example", "domain"}},
			},
			"3": {
				ID:          "3",
				URL:         "http://example.com",
				Title:       "Example Domain",
				Description: "",
				Phrases:     [][]string{{"example", "domain"}},
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

		if info.TotalPhrases != 3 {
			t.Errorf("expected 3 keywords, got %d", info.TotalPhrases)
		}

		if info.SizeOfIndex != len(encoded) {
			t.Errorf("expected %d index size, got %d", len(encoded), info.SizeOfIndex)
		}
	})
}

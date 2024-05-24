package handler_test

import (
	"crawlquery/node/domain"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"
	"encoding/json"

	dumpHandler "crawlquery/node/dump/handler"
	dumpService "crawlquery/node/dump/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPageDump(t *testing.T) {
	t.Run("can do page dump", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		dumpService := dumpService.NewService(pageService)

		_, err := pageService.Create("1", "http://example.com", "hash1")

		if err != nil {
			t.Fatalf("Error saving page: %v", err)
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/hash", nil)

		dumpHandler := dumpHandler.NewHandler(dumpService)

		dumpHandler.Page(ctx)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}

		var slicePages map[string]*domain.Page

		err = json.Unmarshal(w.Body.Bytes(), &slicePages)

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

		if slicePages["1"].Hash != "hash1" {
			t.Fatalf("Expected page Hash to be 'hash1', got '%s'", slicePages["1"].Hash)
		}
	})
}

package handler_test

import (
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

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

		_, err := pageService.Create("1", "http://example.com")

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
		expected := `{"1":{"id":"1","url":"http://example.com","title":"","meta_description":"","keywords":null}}`
		if w.Body.String() != expected {
			t.Fatalf("Expected body to be '%s', got '%s'", expected, w.Body.String())
		}
	})
}

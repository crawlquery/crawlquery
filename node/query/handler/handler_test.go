package handler_test

import (
	"bytes"
	"crawlquery/node/dto"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	queryHandler "crawlquery/node/query/handler"
	queryService "crawlquery/node/query/service"

	"github.com/gin-gonic/gin"
)

func TestQuery(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		queryService := queryService.NewService(nil)

		queryHandler := queryHandler.NewHandler(queryService)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := dto.QueryRequest{
			Query: "SELECT title FROM pages WHERE title LIKE '%example%'",
		}

		reqBody, _ := json.Marshal(req)

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/query", bytes.NewReader(reqBody))

		queryHandler.Query(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var res dto.QueryResponse

		err := json.NewDecoder(w.Body).Decode(&res)

		if err != nil {
			t.Errorf("Error decoding response: %v", err)
		}

		if len(res.Pages) != 1 {
			t.Errorf("Expected 1 result, got %d", len(res.Pages))
		}

		if res.Pages[0].ID != "page1" {
			t.Errorf("Expected page ID page1, got %s", res.Pages[0].ID)
		}

		if res.Pages[0].Title != "Example Page" {
			t.Errorf("Expected title Example Page, got %s", res.Pages[0].Title)
		}

	})
}

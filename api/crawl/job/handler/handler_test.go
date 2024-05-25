package handler_test

import (
	"bytes"
	"crawlquery/api/crawl/job/handler"
	"crawlquery/api/crawl/job/repository/mem"
	"crawlquery/api/crawl/job/service"
	"crawlquery/api/dto"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreate(t *testing.T) {
	t.Run("should create a crawl job", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
		h := handler.NewHandler(svc)

		// given
		a := &dto.CreateCrawlJobRequest{
			URL: "http://example.com",
		}

		req, err := a.ToJSON()

		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/crawl", bytes.NewBuffer(req))

		// when
		h.Create(ctx)

		// then
		if ctx.Writer.Status() != http.StatusCreated {
			t.Errorf("Expected status to be 201, got %d", ctx.Writer.Status())
		}

		var res dto.CreateCrawlJobResponse
		err = json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if res.CrawlJob.URL != a.URL {
			t.Errorf("Expected URL to be %s, got %s", a.URL, res.CrawlJob.URL)
		}
	})

	t.Run("should return 400 if URL is invalid", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
		h := handler.NewHandler(svc)

		// given
		a := &dto.CreateCrawlJobRequest{
			URL: "x123!",
		}

		req, err := a.ToJSON()

		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/crawl", bytes.NewBuffer(req))

		// when
		h.Create(ctx)

		// then
		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}
	})

	t.Run("should return 400 if malformed JSON", func(t *testing.T) {
		// Arrange
		repo := mem.NewRepository()
		svc := service.NewService(repo, nil, nil, nil, nil, nil, testutil.NewTestLogger())
		h := handler.NewHandler(svc)

		// given
		req := []byte(`{"url":`)
		responseWriter := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/crawl", bytes.NewBuffer(req))

		// when
		h.Create(ctx)

		// then
		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}
	})
}

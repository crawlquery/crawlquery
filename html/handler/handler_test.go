package handler_test

import (
	"bytes"
	"crawlquery/html/dto"
	"crawlquery/html/handler"
	"crawlquery/pkg/util"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetPage(t *testing.T) {

	os.MkdirAll("/tmp/cq-html-test", 0755)
	defer os.RemoveAll("/tmp/cq-html-test")

	handler := handler.NewHandler("/tmp/cq-html-test")

	html := []byte("<html><body><h1>Hello, World!</h1></body></html>")

	pageID := "page1"

	if err := os.WriteFile("/tmp/cq-html-test/page1", html, 0644); err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	defer os.Remove("/tmp/cq-html-test/page1")

	// Test valid page
	t.Run("valid page", func(t *testing.T) {
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{
			{Key: "pageID", Value: pageID},
		}

		handler.GetPage(ctx)

		var resp dto.GetPageResponse

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if string(resp.HTML) != string(html) {
			t.Fatalf("Test failed: Expected body %v, got %v", string(html), string(resp.HTML))
		}

		if w.Code != 200 {
			t.Fatalf("Test failed: Expected status 200, got %v", w.Code)
		}
	})

	// Test invalid page
	t.Run("invalid page", func(t *testing.T) {

		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{
			{Key: "pageID", Value: "!!!"},
		}

		handler.GetPage(ctx)

		var errResp dto.ErrorResponse

		if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {

			t.Fatalf("Test failed: %v", err)
		}

		if w.Code != 404 {
			t.Fatalf("Test failed: Expected status 404, got %v", w.Code)
		}

		if errResp.Error != "invalid page ID" {
			t.Fatalf("Test failed: Expected error message 'invalid page ID', got %v", errResp.Error)
		}
	})

	// Test page not found
	t.Run("page not found", func(t *testing.T) {

		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{
			{Key: "pageID", Value: "page2"},
		}

		handler.GetPage(ctx)

		var errResp dto.ErrorResponse

		if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if w.Code != 404 {
			t.Fatalf("Test failed: Expected status 404, got %v", w.Code)
		}

		if errResp.Error != "page not found" {
			t.Fatalf("Test failed: Expected error message 'page not found', got %v", errResp.Error)
		}
	})

	t.Run("page too large", func(t *testing.T) {

		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{
			{Key: "pageID", Value: pageID},
		}

		body := make([]byte, 10000001)

		for i := 0; i < 10000001; i++ {
			body[i] = 'a'
		}

		storeReq := dto.StorePageRequest{
			Hash: string(util.Sha256Hex32([]byte("http://google.com"))),
			HTML: body,
		}

		reqBody, err := json.Marshal(storeReq)

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		ctx.Request = httptest.NewRequest("POST", "/store", bytes.NewBuffer(reqBody))

		handler.StorePage(ctx)

		var errResp dto.ErrorResponse

		if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if w.Code != 400 {
			t.Fatalf("Test failed: Expected status 400, got %v", w.Code)
		}

		if errResp.Error != "invalid request" {
			t.Fatalf("Test failed: Expected error message 'invalid request', got %v", errResp.Error)
		}
	})
}

func TestStorePage(t *testing.T) {

	os.MkdirAll("/tmp/cq-html-test", 0755)
	defer os.RemoveAll("/tmp/cq-html-test")

	handler := handler.NewHandler("/tmp/cq-html-test")

	html := []byte("<html><body><h1>Hello, World!</h1></body></html>")

	// Test valid page
	t.Run("valid page", func(t *testing.T) {
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		body := dto.StorePageRequest{
			Hash: string(util.PageID("http://google.com")),
			HTML: html,
		}

		reqBody, err := json.Marshal(body)

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		ctx.Request = httptest.NewRequest("POST", "/store", bytes.NewBuffer(reqBody))

		handler.StorePage(ctx)

		if w.Code != 201 {
			t.Fatalf("Test failed: Expected status 200, got %v", w.Code)
		}

		data, err := os.ReadFile("/tmp/cq-html-test/" + body.Hash)

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if string(data) != string(html) {
			t.Fatalf("Test failed: Expected body %v, got %v", string(html), string(data))
		}
	})

	// Test invalid page
	t.Run("invalid page", func(t *testing.T) {
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		body := dto.StorePageRequest{
			Hash: "!!!",
			HTML: html,
		}

		reqBody, err := json.Marshal(body)

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		ctx.Request = httptest.NewRequest("POST", "/store", bytes.NewBuffer(reqBody))

		handler.StorePage(ctx)

		if w.Code != 400 {
			t.Fatalf("Test failed: Expected status 400, got %v", w.Code)
		}

		var errResp dto.ErrorResponse

		if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if errResp.Error != "invalid page ID" {
			t.Fatalf("Test failed: Expected error message 'invalid page ID', got %v", errResp.Error)
		}
	})
}

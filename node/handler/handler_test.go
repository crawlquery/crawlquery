package handler_test

import (
	"bytes"
	"crawlquery/api/router"
	"crawlquery/pkg/factory"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchHandler(t *testing.T) {
	r := router.NewRouter()

	req, _ := http.NewRequest("GET", "/search", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expect, err := json.Marshal(gin.H{
		"results": factory.ExampleResults(),
	})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.Bytes()

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	if !bytes.Equal(body, expect) {
		t.Errorf("expected body %v; got %v", expect, body)
	}
}

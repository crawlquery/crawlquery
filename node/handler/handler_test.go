package handler_test

import (
	"crawlquery/node/handler"
	"crawlquery/node/router"
	"crawlquery/node/service"
	"crawlquery/pkg/domain"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/index"
	"crawlquery/pkg/repository/index/mem"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchHandler(t *testing.T) {

	idx := index.NewIndex()
	for _, page := range factory.TenPages() {
		idx.AddPage(page)
	}

	memRepo := mem.NewMemoryRepository()
	memRepo.Save(idx)

	is := service.NewIndexService(
		memRepo,
	)

	handler := handler.NewSearchHandler(is)

	r := router.NewRouter(handler)

	req, _ := http.NewRequest("GET", "/search?q=homepage", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expected, err := json.Marshal(gin.H{
		"results": domain.Result{
			PageID: "doc1",
			Score:  1.0,
			Page: domain.Page{
				ID:    "doc1",
				URL:   "http://google.com",
				Title: "Google",
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	if body != string(expected) {
		t.Errorf("expected body %s; got %s", expected, body)
	}
}

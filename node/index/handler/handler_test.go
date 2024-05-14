package handler_test

import (
	crawlHandler "crawlquery/node/crawl/handler"
	"crawlquery/node/index"
	indexHandler "crawlquery/node/index/handler"
	"crawlquery/node/router"
	"crawlquery/pkg/domain"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	forwardRepository "crawlquery/node/index/forward/repository/mem"
	invertedRepository "crawlquery/node/index/inverted/repository/mem"

	"github.com/gin-gonic/gin"
)

func TestIndexHandler(t *testing.T) {

	idx := index.NewIndex(
		forwardRepository.NewRepository(),
		invertedRepository.NewRepository(),
		testutil.NewTestLogger(),
	)
	for _, page := range factory.TenPages() {
		idx.AddPage(page)
	}

	indexHandler := indexHandler.NewHandler(idx)
	crawlHandler := crawlHandler.NewHandler(nil)

	r := router.NewRouter(indexHandler, crawlHandler)

	req, _ := http.NewRequest("GET", "/search?q=homepage", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expected, err := json.Marshal(gin.H{
		"results": []domain.Result{
			{
				PageID: "1",
				Score:  1.0,
				Page:   factory.HomePage,
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

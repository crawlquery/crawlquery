package handler_test

import (
	"bytes"
	"crawlquery/node/domain"
	indexHandler "crawlquery/node/index/handler"
	sharedDomain "crawlquery/pkg/domain"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	htmlRepo "crawlquery/node/html/repository/mem"
	htmlService "crawlquery/node/html/service"

	indexService "crawlquery/node/index/service"

	peerService "crawlquery/node/peer/service"

	"github.com/gin-gonic/gin"
)

func TestSearch(t *testing.T) {

	t.Run("returns results", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			keywordService,
			peerService,
			testutil.NewTestLogger(),
		)

		for k, dummy := range factory.ThreePages {

			page, err := pageService.Create(k, dummy.URL)

			if err != nil {
				t.Fatalf("error saving page: %v", err)
			}

			err = htmlService.Save(page.ID, []byte(dummy.HTML))

			if err != nil {
				t.Fatalf("error saving html: %v", err)
			}

			err = indexService.Index(page.ID)

			if err != nil {
				t.Fatalf("error indexing page: %v", err)
			}
		}

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search?q=example+website", nil)

		indexHandler.Search(ctx)

		expected, err := json.Marshal(gin.H{
			"results": []sharedDomain.Result{
				{
					PageID: "home1",
					Score:  1.0,
					Page: &sharedDomain.Page{
						ID:              "home1",
						URL:             "https://example.com",
						Title:           "Home",
						MetaDescription: "Welcome to our official website where we offer the latest updates and information.",
					},
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
	})

	t.Run("returns error if query is missing", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			keywordService,
			peerService,
			testutil.NewTestLogger(),
		)

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/search", nil)

		indexHandler.Search(ctx)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status BadRequest; got %v", w.Code)
		}
	})
}

func TestEvent(t *testing.T) {
	t.Run("can handle index event", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo)

		htmlRepo := htmlRepo.NewRepository()
		htmlService := htmlService.NewService(htmlRepo)

		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		peerService := peerService.NewService(nil, keywordService, pageService, nil, testutil.NewTestLogger())

		indexService := indexService.NewService(
			pageService,
			htmlService,
			keywordService,
			peerService,
			testutil.NewTestLogger(),
		)

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		indexEvent := domain.IndexEvent{
			Page: &sharedDomain.Page{
				URL:             "http://example.com",
				ID:              "page1",
				Title:           "Example",
				MetaDescription: "An example page",
			},
			Keywords: map[string]*domain.Posting{
				"keyword1": {
					PageID:    "page1",
					Frequency: 1,
					Positions: []int{1},
				},
			},
		}

		encoded, err := json.Marshal(indexEvent)

		if err != nil {
			t.Fatalf("error encoding index event: %v", err)
		}

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(encoded))

		indexHandler.Event(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		page, err := pageService.Get("page1")

		if err != nil {
			t.Fatalf("error getting page: %v", err)
		}

		if page == nil {
			t.Fatalf("expected page to be found, got nil")
		}

		if page.ID != "page1" {
			t.Fatalf("expected page ID to be page1, got %s", page.ID)
		}

		if page.Title != "Example" {
			t.Fatalf("expected title to be Example, got %s", page.Title)
		}

		if page.MetaDescription != "An example page" {
			t.Fatalf("expected meta description to be An example page, got %s", page.MetaDescription)
		}

		postings, err := keywordService.GetPostings("keyword1")

		if err != nil {
			t.Fatalf("error getting keyword: %v", err)
		}

		if len(postings) != 1 {
			t.Fatalf("expected 1 posting, got %d", len(postings))
		}

		if postings[0].PageID != "page1" {
			t.Fatalf("expected page id to be page1, got %s", postings[0].PageID)
		}

		if postings[0].Frequency != 1 {
			t.Fatalf("expected frequency to be 1, got %d", postings[0].Frequency)
		}
	})
}

func TestHash(t *testing.T) {
	t.Run("returns hash", func(t *testing.T) {
		keywordRepo := keywordRepo.NewRepository()
		keywordService := keywordService.NewService(keywordRepo)

		indexService := indexService.NewService(nil, nil, keywordService, nil, testutil.NewTestLogger())

		postings := map[string]*domain.Posting{
			"token": {
				PageID:    "page1",
				Frequency: 1,
				Positions: []int{1},
			},
		}

		keywordService.SavePostings(postings)

		indexHandler := indexHandler.NewHandler(indexService, testutil.NewTestLogger())

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/hash", nil)

		indexHandler.Hash(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", w.Code)
		}

		var response map[string]string

		err := json.Unmarshal(w.Body.Bytes(), &response)

		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		if response["hash"] == "" {
			t.Fatalf("expected hash to not be empty")
		}
	})
}

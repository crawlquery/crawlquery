package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"crawlquery/node/domain"
	"crawlquery/node/dto"
	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	peerService "crawlquery/node/peer/service"
	repairService "crawlquery/node/repair/service"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	repairHandler "crawlquery/node/repair/handler"

	"github.com/gin-gonic/gin"
)

func TestGetIndexMetas(t *testing.T) {
	t.Run("can get index metas", func(t *testing.T) {
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		now := time.Now()
		oneHourAgo := now.Add(-time.Hour)
		pageRepo.Save("page1", &domain.Page{
			ID:            "page1",
			LastIndexedAt: &now,
		})

		pageRepo.Save("page2", &domain.Page{
			ID:            "page2",
			LastIndexedAt: &oneHourAgo,
		})

		peerService := peerService.NewService(nil, &domain.Peer{
			ID:       "peer1",
			Hostname: "peer1.cluster.com",
			Port:     8080,
		}, nil)

		repairService := repairService.NewService(nil, pageService, nil, peerService, nil)

		handler := repairHandler.NewHandler(repairService)

		var req dto.GetIndexMetasRequest

		req.PageIDs = []string{"page1", "page2"}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Fatal(err)
		}

		ctx.Request = httptest.NewRequest("POST", "/index-metas", bytes.NewBuffer(encoded))

		handler.GetIndexMetas(ctx)

		if w.Code != 200 {
			t.Fatalf("expected status code 200, got %d", w.Code)
		}

		var res dto.GetIndexMetasResponse

		err = json.NewDecoder(w.Body).Decode(&res)

		if err != nil {
			t.Fatal(err)
		}

		if len(res.IndexMetas) != 2 {
			t.Fatalf("expected 2 index metas, got %d", len(res.IndexMetas))
		}

		if res.IndexMetas[0].PageID != "page1" {
			t.Errorf("expected page1, got %s", res.IndexMetas[0].PageID)
		}

		if res.IndexMetas[1].PageID != "page2" {
			t.Errorf("expected page2, got %s", res.IndexMetas[1].PageID)
		}
	})
}

func TestGetPageDumps(t *testing.T) {
	t.Run("can get page dumps", func(t *testing.T) {
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("page1", &domain.Page{
			ID: "page1",
		})

		pageRepo.Save("page2", &domain.Page{
			ID: "page2",
		})

		peerService := peerService.NewService(nil, &domain.Peer{
			ID:       "peer1",
			Hostname: "peer1.cluster.com",
			Port:     8080,
		}, nil)

		keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
		keywordService := keywordService.NewService(keywordOccurrenceRepo)

		repairService := repairService.NewService(nil, pageService, keywordService, peerService, nil)

		handler := repairHandler.NewHandler(repairService)

		var req dto.GetPageDumpsRequest

		req.PageIDs = []string{"page1", "page2"}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Fatal(err)
		}

		ctx.Request = httptest.NewRequest("POST", "/page-dumps", bytes.NewBuffer(encoded))

		handler.GetPageDumps(ctx)

		if w.Code != 200 {
			t.Fatalf("expected status code 200, got %d", w.Code)
		}

		var res dto.GetPageDumpsResponse

		err = json.NewDecoder(w.Body).Decode(&res)

		if err != nil {
			t.Fatal(err)
		}

		if len(res.PageDumps) != 2 {
			t.Fatalf("expected 2 page dumps, got %d", len(res.PageDumps))
		}

		if res.PageDumps[0].PageID != "page1" {
			t.Errorf("expected page1, got %s", res.PageDumps[0].PageID)
		}

		if res.PageDumps[1].PageID != "page2" {
			t.Errorf("expected page2, got %s", res.PageDumps[1].PageID)
		}
	})
}

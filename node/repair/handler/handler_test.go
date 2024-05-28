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

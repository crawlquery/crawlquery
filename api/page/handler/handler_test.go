package handler_test

import (
	"bytes"
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"encoding/json"
	"net/http/httptest"
	"testing"

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlLogRepo "crawlquery/api/crawl/log/repository/mem"

	crawlService "crawlquery/api/crawl/service"

	eventService "crawlquery/api/event/service"

	pageHandler "crawlquery/api/page/handler"

	"github.com/gin-gonic/gin"
)

func setup() (*pageRepo.Repository, *pageService.Service) {
	eventService := eventService.NewService()
	pageRepo := pageRepo.NewRepository()
	shardRepo := shardRepo.NewRepository()

	shardRepo.Create(&domain.Shard{
		ID: 0,
	})
	crawlJobRepo := crawlJobRepo.NewRepository()
	crawlLogRepo := crawlLogRepo.NewRepository()
	crawlService := crawlService.NewService(
		crawlService.WithCrawlJobRepo(crawlJobRepo),
		crawlService.WithCrawlLogRepo(crawlLogRepo),
		crawlService.WithLogger(testutil.NewTestLogger()),
	)
	shardService := shardService.NewService(
		shardService.WithRepo(shardRepo),
		shardService.WithLogger(testutil.NewTestLogger()),
	)
	pageService := pageService.NewService(
		pageService.WithPageRepo(pageRepo),
		pageService.WithShardService(shardService),
		pageService.WithCrawlService(crawlService),
		pageService.WithLogger(testutil.NewTestLogger()),
		pageService.WithEventService(eventService),
	)

	return pageRepo, pageService
}

func TestCreate(t *testing.T) {
	t.Run("can create a page", func(t *testing.T) {
		// setup
		pageRepo, pageService := setup()

		pageHandler := pageHandler.NewHandler(pageService)

		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		var req dto.CreatePageRequest
		req.URL = "http://example.com"

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Fatalf("error encoding request: %v", err)
		}

		ctx.Request = httptest.NewRequest("POST", "/pages", bytes.NewBuffer(encoded))

		pageHandler.Create(ctx)

		if w.Code != 201 {
			t.Errorf("got status code %d, want 201", w.Code)
		}

		// test
		url := domain.URL("http://example.com")
		pageID := util.PageID(url)

		page, err := pageRepo.Get(pageID)

		if err != nil {
			t.Fatalf("error getting page: %v", err)
		}

		if page.ID != pageID {
			t.Errorf("got page ID %s, want %s", page.ID, pageID)
		}

		if page.URL != url {
			t.Errorf("got page URL %s, want %s", page.URL, url)
		}
	})
}

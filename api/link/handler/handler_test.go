package handler_test

import (
	"bytes"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"encoding/json"
	"net/http/httptest"

	"crawlquery/api/dto"
	linkHandler "crawlquery/api/link/handler"
	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"

	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlJobService "crawlquery/api/crawl/job/service"

	crawlRestrictionRepo "crawlquery/api/crawl/restriction/repository/mem"
	crawlRestrictionService "crawlquery/api/crawl/restriction/service"

	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreate(t *testing.T) {
	t.Run("can create a link", func(t *testing.T) {
		// Arrange
		linkRepo := linkRepo.NewRepository()

		crawlRestrictionRepo := crawlRestrictionRepo.NewRepository()
		crawlRestrictionService := crawlRestrictionService.NewService(crawlRestrictionRepo, testutil.NewTestLogger())

		crawlJobRepo := crawlJobRepo.NewRepository()
		crawlJobService := crawlJobService.NewService(crawlJobRepo, nil, nil, crawlRestrictionService, nil, nil, testutil.NewTestLogger())

		linkService := linkService.NewService(linkRepo, crawlJobService, testutil.NewTestLogger())

		linkHandler := linkHandler.NewHandler(linkService, testutil.NewTestLogger())

		src := "https://cancreatealink.com"
		dst := "https://cancreatealink.com/about"

		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		reqBody := &dto.CreateLinkRequest{
			Src: src,
			Dst: dst,
		}

		encoded, err := json.Marshal(reqBody)

		if err != nil {
			t.Fatalf("error encoding request body: %v", err)
		}

		ctx.Request = httptest.NewRequest("POST", "/links", bytes.NewReader(encoded))

		linkHandler.Create(ctx)

		repoCheck, _ := linkRepo.GetAllBySrcID(util.PageID(src))

		if len(repoCheck) != 1 {
			t.Errorf("Expected 1 link, got %d", len(repoCheck))
		}

		if ctx.Writer.Status() != 201 {
			t.Errorf("expected status Created; got %v", ctx.Writer.Status())
		}

		if w.Body.String() != "" {
			t.Errorf("expected empty body; got %v", w.Body.String())
		}

		if repoCheck[0].SrcID != util.PageID(src) {
			t.Errorf("Expected srcID %s, got %s", util.PageID(src), repoCheck[0].SrcID)
		}

		if repoCheck[0].DstID != util.PageID(dst) {
			t.Errorf("Expected dstID %s, got %s", util.PageID(dst), repoCheck[0].DstID)
		}
	})
}

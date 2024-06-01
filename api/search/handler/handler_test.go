package handler_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/search/handler"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	nodeDomain "crawlquery/node/domain"
	nodeDto "crawlquery/pkg/dto"

	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"

	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"

	pageRankRepo "crawlquery/api/pagerank/repository/mem"
	pageRankService "crawlquery/api/pagerank/service"

	searchService "crawlquery/api/search/service"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
)

func setupServices() (*nodeRepo.Repository, *nodeService.Service, *linkService.Service, *pageRankService.Service, *searchService.Service) {
	nodeRepo := nodeRepo.NewRepository()
	nodeService := nodeService.NewService(
		nodeService.WithNodeRepo(nodeRepo),
		nodeService.WithLogger(testutil.NewTestLogger()),
		nodeService.WithRandSeed(time.Now().Unix()),
	)
	linkRepo := linkRepo.NewRepository()
	linkService := linkService.NewService(
		linkService.WithLinkRepo(linkRepo),
		linkService.WithLogger(testutil.NewTestLogger()),
	)
	pageRankRepo := pageRankRepo.NewRepository()
	pageRankService := pageRankService.NewService(linkService, pageRankRepo, testutil.NewTestLogger())
	searchService := searchService.NewService(nodeService, pageRankService, testutil.NewTestLogger())

	return nodeRepo, nodeService, linkService, pageRankService, searchService
}
func TestSearch(t *testing.T) {
	t.Run("should return results", func(t *testing.T) {

		nodeRepo, _, _, _, searchService := setupServices()

		nodeRepo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		nodeRepo.Create(&domain.Node{
			ID:        "node2",
			ShardID:   1,
			Hostname:  "node2.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Get("/search").
			MatchParam("q", "term").
			Reply(200).
			JSON(nodeDto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page1",
						Score:  0.5,
						Page: nodeDomain.ResultPage{
							ID:    "page1",
							URL:   "http://google.com",
							Title: "Google",
						},
					},
				},
			})

		gock.New("http://node2.cluster.com:8080").
			Get("/search").
			MatchParam("q", "term").
			Reply(200).
			JSON(nodeDto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page2",
						Score:  0.6,
						Page: nodeDomain.ResultPage{
							ID:    "page2",
							URL:   "http://facebook.com",
							Title: "Facebook",
						},
					},
				},
			})

		handler := handler.NewHandler(searchService)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("GET", "/search?q=term", nil)

		// when
		handler.Search(ctx)

		// then
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status to be 200, got %d", ctx.Writer.Status())
		}

		var res dto.SearchResponse
		err := json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if len(res.Results) != 2 {
			t.Errorf("Expected 1 result, got %d", len(res.Results))
		}
	})
}

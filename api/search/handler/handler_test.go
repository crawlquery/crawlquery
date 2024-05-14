package handler_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/node/repository/mem"
	"crawlquery/api/search/handler"
	"crawlquery/api/search/service"
	"crawlquery/pkg/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sharedDomain "crawlquery/pkg/domain"
	nodeDto "crawlquery/pkg/dto"

	nodeService "crawlquery/api/node/service"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
)

func TestSearch(t *testing.T) {
	t.Run("should return results", func(t *testing.T) {

		repo := mem.NewRepository()

		repo.Create(&domain.Node{
			ID:        "node1",
			ShardID:   0,
			Hostname:  "node1.cluster.com",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		repo.Create(&domain.Node{
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
				Results: []sharedDomain.Result{
					{
						PageID: "page1",
						Score:  0.5,
						Page: &sharedDomain.Page{
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
				Results: []sharedDomain.Result{
					{
						PageID: "page2",
						Score:  0.6,
						Page: &sharedDomain.Page{
							ID:    "page2",
							URL:   "http://facebook.com",
							Title: "Facebook",
						},
					},
				},
			})

		nodeService := nodeService.NewService(repo, nil, nil, testutil.NewTestLogger())

		svc := service.NewService(nodeService, testutil.NewTestLogger())
		handler := handler.NewHandler(svc)

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

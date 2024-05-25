package service_test

import (
	"crawlquery/api/domain"
	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"
	"crawlquery/api/search/service"
	nodeDomain "crawlquery/node/domain"
	"crawlquery/pkg/dto"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestSearch(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()

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

		nodeService := nodeService.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		defer gock.Off()

		gock.New("http://node1.cluster.com:8080").
			Get("/search").
			MatchParam("q", "term").
			Reply(200).
			JSON(dto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page1",
						Score:  0.5,
						Page: &nodeDomain.ResultPage{
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
			JSON(dto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page2",
						Score:  0.6,
						Page: &nodeDomain.ResultPage{
							ID:    "page2",
							URL:   "http://facebook.com",
							Title: "Facebook",
						},
					},
				},
			})

		searchService := service.NewService(nodeService, testutil.NewTestLogger())

		results, err := searchService.Search("term")
		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %v", len(results))
		}

		if results[0].PageID != "page1" && results[1].PageID != "page1" {
			t.Errorf("Expected page ID to be page1, got %s and %s", results[0].PageID, results[1].PageID)
		}
	})

	t.Run("cleans term", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()

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
			MatchParam("q", "term hello").
			Reply(200).
			JSON(dto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page1",
						Score:  0.5,
						Page: &nodeDomain.ResultPage{
							ID:    "page1",
							URL:   "http://google.com",
							Title: "Google",
						},
					},
				},
			})

		gock.New("http://node2.cluster.com:8080").
			Get("/search").
			MatchParam("q", "term hello").
			Reply(200).
			JSON(dto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page1",
						Score:  0.5,
						Page: &nodeDomain.ResultPage{
							ID:    "page1",
							URL:   "http://google.com",
							Title: "Google",
						},
					},
				},
			})

		nodeService := nodeService.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())
		searchService := service.NewService(nodeService, testutil.NewTestLogger())

		results, err := searchService.Search("   term      hello   ")
		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %v", len(results))
		}
	})
	t.Run("removes duplicate results", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()

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
			JSON(dto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page1",
						Score:  0.5,
						Page: &nodeDomain.ResultPage{
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
			JSON(dto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page1",
						Score:  0.5,
						Page: &nodeDomain.ResultPage{
							ID:    "page1",
							URL:   "http://google.com",
							Title: "Google",
						},
					},
				},
			})

		nodeService := nodeService.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())
		searchService := service.NewService(nodeService, testutil.NewTestLogger())

		results, err := searchService.Search("term")
		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %v", len(results))
		}
	})
}

package service_test

import (
	"crawlquery/api/domain"
	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"
	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"
	nodeDomain "crawlquery/node/domain"
	"math"

	pageRankRepo "crawlquery/api/pagerank/repository/mem"
	pageRankService "crawlquery/api/pagerank/service"
	searchService "crawlquery/api/search/service"

	"crawlquery/pkg/dto"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func setupServices() (*nodeRepo.Repository, *nodeService.Service, *linkService.Service, *linkRepo.Repository, *pageRankService.Service, *searchService.Service) {
	nodeRepo := nodeRepo.NewRepository()
	nodeService := nodeService.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())
	linkRepo := linkRepo.NewRepository()
	linkService := linkService.NewService(linkRepo, nil, testutil.NewTestLogger())
	pageRankRepo := pageRankRepo.NewRepository()
	pageRankService := pageRankService.NewService(linkService, pageRankRepo, testutil.NewTestLogger())
	searchService := searchService.NewService(nodeService, pageRankService, testutil.NewTestLogger())

	return nodeRepo, nodeService, linkService, linkRepo, pageRankService, searchService
}

func TestSearch(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		nodeRepo, _, _, _, _, searchService := setupServices()

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
			JSON(dto.NodeSearchResponse{
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
		nodeRepo, _, _, _, _, searchService := setupServices()

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
			MatchParam("q", "term hello").
			Reply(200).
			JSON(dto.NodeSearchResponse{
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

		results, err := searchService.Search("   term      hello   ")
		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %v", len(results))
		}
	})
	t.Run("removes duplicate results", func(t *testing.T) {
		nodeRepo, _, _, _, _, searchService := setupServices()

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
			JSON(dto.NodeSearchResponse{
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

		results, err := searchService.Search("term")
		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %v", len(results))
		}
	})

	t.Run("applies page rank", func(t *testing.T) {
		nodeRepo, _, _, linkRepo, _, searchService := setupServices()

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
			JSON(dto.NodeSearchResponse{
				Results: []nodeDomain.Result{
					{
						PageID: "page2",
						Page: nodeDomain.ResultPage{
							ID:    "page2",
							URL:   "http://facebook.com",
							Title: "Facebook",
						},
					},
				},
			})

		linkRepo.Create(&domain.Link{
			SrcID: "page3",
			DstID: "page1",
		})

		linkRepo.Create(&domain.Link{
			SrcID: "page3",
			DstID: "page1",
		})

		results, err := searchService.Search("term")
		if err != nil {
			t.Fatalf("Error searching: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %v", len(results))
		}

		if results[0].PageID != "page2" && results[1].PageID != "page2" {
			t.Errorf("Expected page ID to be page2, got %s and %s", results[0].PageID, results[1].PageID)
		}

		expectedRank := 0.1387

		for _, result := range results {
			if result.PageID == "page1" {
				if math.Abs(result.PageRank-expectedRank) > 0.01 {
					t.Errorf("Unexpected PageRank for page1. Expected %f, got %f", expectedRank, result.PageRank)
				}
			}
		}
	})
}

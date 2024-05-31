package service_test

import (
	"crawlquery/api/domain"
	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"
	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"
	nodeDomain "crawlquery/node/domain"

	pageRankRepo "crawlquery/api/pagerank/repository/mem"
	pageRankService "crawlquery/api/pagerank/service"
	searchService "crawlquery/api/search/service"

	"crawlquery/pkg/dto"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func setupServices() (*nodeRepo.Repository, *nodeService.Service, *linkService.Service, *linkRepo.Repository, *pageRankRepo.Repository, *pageRankService.Service, *searchService.Service) {
	nodeRepo := nodeRepo.NewRepository()
	nodeService := nodeService.NewService(
		nodeService.WithNodeRepo(nodeRepo),
		nodeService.WithLogger(testutil.NewTestLogger()),
	)
	linkRepo := linkRepo.NewRepository()
	linkService := linkService.NewService(
		linkService.WithLinkRepo(linkRepo),
		linkService.WithLogger(testutil.NewTestLogger()),
	)
	pageRankRepo := pageRankRepo.NewRepository()
	pageRankService := pageRankService.NewService(linkService, pageRankRepo, testutil.NewTestLogger())
	searchService := searchService.NewService(nodeService, pageRankService, testutil.NewTestLogger())

	return nodeRepo, nodeService, linkService, linkRepo, pageRankRepo, pageRankService, searchService
}

func TestSearch(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		nodeRepo, _, _, _, _, _, searchService := setupServices()

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
		nodeRepo, _, _, _, _, _, searchService := setupServices()

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
		nodeRepo, _, _, _, _, _, searchService := setupServices()

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
		nodeRepo, _, _, _, pageRankRepo, _, searchService := setupServices()

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

		pageRankRepo.Update("page1", 0.5, time.Now())

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

		for _, result := range results {
			if result.PageID == "page1" {
				if result.PageRank != 0.5 {
					t.Errorf("Unexpected PageRank for page1. Expected %f, got %f", 0.5, result.PageRank)
				}
			}
		}
	})
}

package service_test

import (
	"crawlquery/api/service"
	"crawlquery/pkg/domain"
	"crawlquery/pkg/dto"
	nodeMemRepo "crawlquery/pkg/repository/node/mem"
	"testing"

	"github.com/h2non/gock"
)

func TestSearch(t *testing.T) {
	nodeRepo := nodeMemRepo.NewMemoryRepository()

	nodeRepo.CreateOrUpdate(&domain.Node{
		ID:       "node1",
		ShardID:  0,
		Hostname: "node1.cluster.com",
		Port:     8080,
	})

	nodeRepo.CreateOrUpdate(&domain.Node{
		ID:       "node2",
		ShardID:  1,
		Hostname: "node2.cluster.com",
		Port:     8080,
	})

	nodeService := service.NewNodeService(nodeRepo)

	defer gock.Off()

	gock.New("http://node1.cluster.com:8080").
		Get("/search").
		MatchParam("q", "term").
		Reply(200).
		JSON(dto.NodeSearchResponse{
			Results: []domain.Result{
				{
					PageID: "page1",
					Score:  0.5,
					Page: domain.Page{
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
			Results: []domain.Result{
				{
					PageID: "page2",
					Score:  0.6,
					Page: domain.Page{
						ID:    "page2",
						URL:   "http://facebook.com",
						Title: "Facebook",
					},
				},
			},
		})

	searchService := service.NewSearchService(nodeService)

	results, err := searchService.Search("term")
	if err != nil {
		t.Fatalf("Error searching: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %v", len(results))
	}

	if results[0].PageID != "page2" {
		t.Errorf("Expected first result to be page2, got %v", results[0].PageID)
	}

	if results[1].PageID != "page1" {
		t.Errorf("Expected second result to be page1, got %v", results[1].PageID)
	}
}

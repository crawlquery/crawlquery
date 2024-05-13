package service_test

import (
	"crawlquery/api/domain"
	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"
	"crawlquery/api/search/service"
	sharedDomain "crawlquery/pkg/domain"
	"crawlquery/pkg/dto"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestSearch(t *testing.T) {
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
		JSON(dto.NodeSearchResponse{
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

	searchService := service.NewService(nodeService, testutil.NewTestLogger())

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

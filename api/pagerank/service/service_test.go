package service_test

import (
	"crawlquery/api/domain"
	linksRepo "crawlquery/api/link/repository/mem"
	linksService "crawlquery/api/link/service"
	pageRankService "crawlquery/api/pagerank/service"
	"math"
	"testing"

	nodeDomain "crawlquery/node/domain"
)

func TestApplyPageRankToResults(t *testing.T) {
	// Create new repository and services
	linksRepo := linksRepo.NewRepository()
	linksService := linksService.NewService(linksRepo, nil, nil)
	pageRankService := pageRankService.NewService(linksService, nil)

	// Create some links for each page
	links := []*domain.Link{
		{SrcID: "A", DstID: "B"},
		{SrcID: "A", DstID: "C"},
		{SrcID: "B", DstID: "C"},
		{SrcID: "C", DstID: "A"},
		{SrcID: "D", DstID: "C"},
		{SrcID: "D", DstID: "E"},
	}

	// Add the links to the repository
	for _, link := range links {
		if err := linksRepo.Create(link); err != nil {
			t.Fatalf("Failed to create link: %v", err)
		}
	}

	// Create a list of results
	results := []nodeDomain.Result{
		{PageID: "A"},
		{PageID: "B"},
		{PageID: "C"},
		{PageID: "D"},
		{PageID: "E"},
	}

	var err error
	// Apply PageRank to the results
	results, err = pageRankService.ApplyPageRankToResults(results)
	if err != nil {
		t.Fatalf("Failed to apply PageRank to results: %v", err)
	}

	// Define expected ranks with a tolerance for comparison
	expectedRanks := map[string]float64{
		"A": 0.26,
		"B": 0.14,
		"C": 0.27,
		"D": 0.03,
		"E": 0.04,
	}

	tolerance := 0.01 // Adjust the tolerance as needed

	for _, result := range results {
		expectedRank, ok := expectedRanks[result.PageID]
		if !ok {
			t.Fatalf("No expected rank found for PageID %s", result.PageID)
		}

		if math.Abs(result.PageRank-expectedRank) > tolerance {
			t.Errorf("Unexpected PageRank for %s. Expected %f, got %f", result.PageID, expectedRank, result.PageRank)
		}
	}
}

func TestCalculatePageRank(t *testing.T) {
	// Create new repository and services
	linksRepo := linksRepo.NewRepository()
	linksService := linksService.NewService(linksRepo, nil, nil)
	pageRankService := pageRankService.NewService(linksService, nil)

	// Create some links for each page
	links := []*domain.Link{
		{SrcID: "A", DstID: "B"},
		{SrcID: "A", DstID: "C"},
		{SrcID: "B", DstID: "C"},
		{SrcID: "C", DstID: "A"},
		{SrcID: "D", DstID: "C"},
		{SrcID: "D", DstID: "E"},
	}

	// Add the links to the repository
	for _, link := range links {
		if err := linksRepo.Create(link); err != nil {
			t.Fatalf("Failed to create link: %v", err)
		}
	}

	expectedRanks := map[string]float64{
		"A": 0.26,
		"B": 0.14,
		"C": 0.27,
		"D": 0.03,
		"E": 0.04,
	}

	tolerance := 0.01 // Adjust the tolerance as needed

	for id, expectedRank := range expectedRanks {
		rank, err := pageRankService.CalculatePageRank(id)
		if err != nil {
			t.Fatalf("Failed to calculate PageRank for %s: %v", id, err)
		}

		if math.Abs(rank-expectedRank) > tolerance {
			t.Errorf("Unexpected PageRank for %s. Expected %f, got %f", id, expectedRank, rank)
		}
	}
}

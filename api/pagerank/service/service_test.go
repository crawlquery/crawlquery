package service_test

import (
	"crawlquery/api/domain"
	linksRepo "crawlquery/api/link/repository/mem"
	linksService "crawlquery/api/link/service"
	pageRankRepo "crawlquery/api/pagerank/repository/mem"
	pageRankService "crawlquery/api/pagerank/service"
	"crawlquery/pkg/testutil"
	"math"
	"testing"
	"time"
)

func TestUpdatePageRanks(t *testing.T) {
	// Create new repository and services
	linksRepo := linksRepo.NewRepository()
	linksService := linksService.NewService(linksRepo, nil, nil)
	pageRankRepo := pageRankRepo.NewRepository()
	pageRankService := pageRankService.NewService(linksService, pageRankRepo, testutil.NewTestLogger())

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

	// Apply PageRank to the results
	err := pageRankService.UpdatePageRanks()
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

	for pageID, expected := range expectedRanks {
		actual, err := pageRankService.GetPageRank(pageID)
		if err != nil {
			t.Fatalf("Failed to get PageRank for %s: %v", pageID, err)
		}

		if math.Abs(actual-expected) > tolerance {
			t.Errorf("Unexpected PageRank for %s. Expected %f, got %f", pageID, expected, actual)
		}
	}
}

func TestGetPageRank(t *testing.T) {
	// Create new repository and services
	linksRepo := linksRepo.NewRepository()
	linksService := linksService.NewService(linksRepo, nil, nil)
	pageRankRepo := pageRankRepo.NewRepository()
	pageRankService := pageRankService.NewService(linksService, pageRankRepo, testutil.NewTestLogger())

	expectedRanks := map[string]float64{
		"A": 0.26,
		"B": 0.14,
		"C": 0.27,
		"D": 0.03,
		"E": 0.04,
	}

	for pageID, expectedRank := range expectedRanks {
		pageRankRepo.Update(pageID, expectedRank, time.Now())
	}

	tolerance := 0.01 // Adjust the tolerance as needed

	for id, expectedRank := range expectedRanks {
		rank, err := pageRankService.GetPageRank(id)
		if err != nil {
			t.Fatalf("Failed to calculate PageRank for %s: %v", id, err)
		}

		if math.Abs(rank-expectedRank) > tolerance {
			t.Errorf("Unexpected PageRank for %s. Expected %f, got %f", id, expectedRank, rank)
		}
	}
}

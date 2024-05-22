package service_test

import (
	"crawlquery/node/domain"
	keywordRepo "crawlquery/node/keyword/repository/mem"
	"crawlquery/node/keyword/service"
	"testing"
)

func TestService(t *testing.T) {
	keywordPostings := map[string]*domain.Posting{
		"test1": {
			PageID:    "page1",
			Frequency: 1,
			Positions: []int{0},
		},
		"test2": {
			PageID:    "page1",
			Frequency: 2,
			Positions: []int{1, 2},
		},
	}

	repo := keywordRepo.NewRepository()

	s := service.NewService(repo)

	err := s.SavePostings(keywordPostings)

	if err != nil {
		t.Fatalf("Error saving postings: %v", err)
	}

	err = s.SavePostings(map[string]*domain.Posting{
		"test2": {
			PageID:    "page2",
			Frequency: 1,
			Positions: []int{0},
		},
	})

	if err != nil {
		t.Fatalf("Error saving postings: %v", err)
	}

	postings, err := repo.GetPostings("test1")

	if err != nil {
		t.Fatalf("Error getting postings: %v", err)
	}

	if len(postings) != 1 {
		t.Fatalf("Expected 1 posting, got %d", len(postings))
	}

	if postings[0].PageID != "page1" {
		t.Fatalf("Expected page id to be page1, got %s", postings[0].PageID)
	}

	if postings[0].Frequency != 1 {
		t.Fatalf("Expected frequency to be 1, got %d", postings[0].Frequency)
	}

	postings, err = repo.GetPostings("test2")

	if err != nil {
		t.Fatalf("Error getting postings: %v", err)
	}

	if len(postings) != 2 {
		t.Fatalf("Expected 2 postings, got %d", len(postings))
	}

	if postings[0].PageID != "page1" {
		t.Fatalf("Expected page id to be page1, got %s", postings[0].PageID)
	}

	if postings[0].Frequency != 2 {
		t.Fatalf("Expected frequency to be 2, got %d", postings[0].Frequency)
	}

	if postings[1].PageID != "page2" {
		t.Fatalf("Expected page id to be page2, got %s", postings[1].PageID)
	}

	if postings[1].Frequency != 1 {
		t.Fatalf("Expected frequency to be 1, got %d", postings[1].Frequency)
	}

	hash, err := repo.GetHash("test1")

	if err != nil {
		t.Fatalf("Error getting hash: %v", err)
	}

	if hash == "" {
		t.Fatalf("Expected hash to not be empty")
	}
}

func TestFuzzySearch(t *testing.T) {
	keywordPostings := map[string]*domain.Posting{
		"test1": {
			PageID:    "page1",
			Frequency: 1,
			Positions: []int{0},
		},
		"test2": {
			PageID:    "page1",
			Frequency: 2,
			Positions: []int{1, 2},
		},
		"notinresults": {
			PageID:    "page2",
			Frequency: 1,
			Positions: []int{0},
		},
	}

	repo := keywordRepo.NewRepository()

	s := service.NewService(repo)

	err := s.SavePostings(keywordPostings)

	if err != nil {
		t.Fatalf("Error saving postings: %v", err)
	}

	results, err := s.FuzzySearch("te")

	if err != nil {
		t.Fatalf("Error fuzzy searching: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0] != "test1" && results[1] != "test1" {
		t.Fatalf("Expected test1 to be in results")
	}

	if results[0] != "test2" && results[1] != "test2" {
		t.Fatalf("Expected test2 to be in results")
	}
}

func TestRemovePostingsByPageID(t *testing.T) {
	keywordPostings := map[string]*domain.Posting{
		"test1": {
			PageID:    "page1",
			Frequency: 1,
			Positions: []int{0},
		},
		"test2": {
			PageID:    "page1",
			Frequency: 2,
			Positions: []int{1, 2},
		},
		"notinresults": {
			PageID:    "page2",
			Frequency: 1,
			Positions: []int{0},
		},
	}

	repo := keywordRepo.NewRepository()

	s := service.NewService(repo)

	err := s.SavePostings(keywordPostings)

	if err != nil {
		t.Fatalf("Error saving postings: %v", err)
	}

	err = s.RemovePostingsByPageID("page1")

	if err != nil {
		t.Fatalf("Error removing postings: %v", err)
	}

	postings, err := repo.GetPostings("test1")

	if err != nil {
		t.Fatalf("Error getting postings: %v", err)
	}

	if len(postings) != 0 {
		t.Fatalf("Expected 0 postings, got %d", len(postings))
	}

	postings, err = repo.GetPostings("test2")

	if err != nil {
		t.Fatalf("Error getting postings: %v", err)
	}

	if len(postings) != 0 {
		t.Fatalf("Expected 0 postings, got %d", len(postings))
	}

	postings, err = repo.GetPostings("notinresults")

	if err != nil {
		t.Fatalf("Error getting postings: %v", err)
	}

	if len(postings) != 1 {
		t.Fatalf("Expected 1 posting, got %d", len(postings))
	}
}

func TestHash(t *testing.T) {
	keywordPostings := map[string]*domain.Posting{
		"test1": {
			PageID:    "page1",
			Frequency: 1,
			Positions: []int{0},
		},
		"test2": {
			PageID:    "page1",
			Frequency: 2,
			Positions: []int{1, 2},
		},
	}

	repo := keywordRepo.NewRepository()

	s := service.NewService(repo)

	err := s.SavePostings(keywordPostings)

	if err != nil {
		t.Fatalf("Error saving postings: %v", err)
	}

	hash, err := repo.GetHash("test1")

	if err != nil {
		t.Fatalf("Error getting hash: %v", err)
	}

	if hash == "" {
		t.Fatalf("Expected hash to not be empty")
	}

	err = s.RemovePostingsByPageID("page1")

	if err != nil {
		t.Fatalf("Error removing postings: %v", err)
	}

	hash, err = repo.GetHash("test1")

	if err != nil {
		t.Fatalf("Error getting hash: %v", err)
	}

	if hash == "" {
		t.Fatalf("Expected hash to not be empty")
	}

	indexHash, err := s.Hash()

	if err != nil {
		t.Fatalf("Error getting hash: %v", err)
	}

	if indexHash == "" {
		t.Fatalf("Expected index hash to not be empty")
	}

	if indexHash != "998b02e885b75a8b874fd86c2bd8caf539590e65444569f0598997ff0fccf4f9" {
		t.Fatalf("Expected hash to be 998b02e885b75a8b874fd86c2bd8caf539590e65444569f0598997ff0fccf4f9, got %s", indexHash)
	}
}

func TestJSON(t *testing.T) {
	keywordPostings := map[string]*domain.Posting{
		"test1": {
			PageID:    "page1",
			Frequency: 1,
			Positions: []int{0},
		},
		"test2": {
			PageID:    "page1",
			Frequency: 2,
			Positions: []int{1, 2},
		},
	}

	repo := keywordRepo.NewRepository()

	s := service.NewService(repo)

	err := s.SavePostings(keywordPostings)

	if err != nil {
		t.Fatalf("Error saving postings: %v", err)
	}

	json, err := s.JSON()

	if err != nil {
		t.Fatalf("Error getting json: %v", err)
	}

	expected := `{"test1":[{"page_id":"page1","frequency":1,"positions":[0]}],"test2":[{"page_id":"page1","frequency":2,"positions":[1,2]}]}`

	if string(json) != expected {
		t.Fatalf("Expected json to be %s, got %s", expected, json)
	}
}

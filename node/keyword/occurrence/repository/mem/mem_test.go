package mem_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/keyword/occurrence/repository/mem"
	"fmt"
	"reflect"
	"testing"
)

func TestGetOccurrences(t *testing.T) {
	repo := mem.NewRepository()
	keyword := domain.Keyword("example")

	// Add occurrences to the repository
	occurrences := []domain.KeywordOccurrence{
		{PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
		{PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}
	for _, occ := range occurrences {
		err := repo.Add(keyword, occ)
		if err != nil {
			t.Fatalf("Error adding occurrence: %v", err)
		}
	}

	// Test GetOccurrences
	gotOccurrences, err := repo.GetAll(keyword)
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	if !reflect.DeepEqual(gotOccurrences, occurrences) {
		t.Errorf("Expected occurrences %v, got %v", occurrences, gotOccurrences)
	}

	// Test GetOccurrences for a non-existing keyword
	_, err = repo.GetAll(domain.Keyword("nonexistent"))
	if err != domain.ErrKeywordNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrKeywordNotFound, err)
	}
}

func TestAddOccurence(t *testing.T) {
	repo := mem.NewRepository()
	keyword := domain.Keyword("example")
	occurrence := domain.KeywordOccurrence{PageID: "page1", Frequency: 1, Positions: []int{1}}

	err := repo.Add(keyword, occurrence)
	if err != nil {
		t.Fatalf("Error adding occurrence: %v", err)
	}

	gotOccurrences, err := repo.GetAll(keyword)
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	if len(gotOccurrences) != 1 {
		t.Fatalf("Expected 1 occurrence, got %d", len(gotOccurrences))
	}

	if !reflect.DeepEqual(gotOccurrences[0], occurrence) {
		t.Errorf("Expected occurrence %v, got %v", occurrence, gotOccurrences[0])
	}
}

func TestGetOccurrencesForPageID(t *testing.T) {
	repo := mem.NewRepository()
	keyword := domain.Keyword("example")

	// Add occurrences to the repository
	occurrences := []domain.KeywordOccurrence{
		{PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
		{PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}
	for _, occ := range occurrences {
		err := repo.Add(keyword, occ)
		if err != nil {
			t.Fatalf("Error adding occurrence: %v", err)
		}
	}

	// Get occurrences for page1
	gotOccurrences, err := repo.GetForPageID("page1")
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	expectedOccurrences := map[domain.Keyword]domain.KeywordOccurrence{
		keyword: {PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
	}

	if !reflect.DeepEqual(gotOccurrences, expectedOccurrences) {
		t.Errorf("Expected occurrences %v, got %v", expectedOccurrences, gotOccurrences)
	}

	// Get occurrences for page2
	gotOccurrences, err = repo.GetForPageID("page2")
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	expectedOccurrences = map[domain.Keyword]domain.KeywordOccurrence{
		keyword: {PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}

	if !reflect.DeepEqual(gotOccurrences, expectedOccurrences) {
		t.Errorf("Expected occurrences %v, got %v", expectedOccurrences, gotOccurrences)
	}

	// Get occurrences for a non-existing page
	gotOccurrences, err = repo.GetForPageID("nonexistent")
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	if len(gotOccurrences) != 0 {
		t.Errorf("Expected no occurrences, got %v", gotOccurrences)
	}
}

func TestRemoveOccurencesForPageID(t *testing.T) {
	repo := mem.NewRepository()
	keyword := domain.Keyword("example")

	// Add occurrences to the repository
	occurrences := []domain.KeywordOccurrence{
		{PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
		{PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}
	for _, occ := range occurrences {
		err := repo.Add(keyword, occ)
		if err != nil {
			t.Fatalf("Error adding occurrence: %v", err)
		}
	}

	// Remove occurrences for page1
	err := repo.RemoveForPageID("page1")
	if err != nil {
		t.Fatalf("Error removing occurrences: %v", err)
	}

	gotOccurrences, err := repo.GetAll(keyword)
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	expectedOccurrences := []domain.KeywordOccurrence{
		{PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}

	if !reflect.DeepEqual(gotOccurrences, expectedOccurrences) {
		t.Errorf("Expected occurrences %v, got %v", expectedOccurrences, gotOccurrences)
	}

	// Remove occurrences for page2
	err = repo.RemoveForPageID("page2")
	if err != nil {
		t.Fatalf("Error removing occurrences: %v", err)
	}

	occ, err := repo.GetAll(keyword)
	if err != domain.ErrKeywordNotFound {
		t.Errorf("Expected no occurrences, got %v", occ)
	}
}

func TestCount(t *testing.T) {
	repo := mem.NewRepository()

	keywords := []domain.Keyword{"example1", "example2", "example3", "example4", "example5"}

	for i, keyword := range keywords {
		for j := 0; j < i+1; j++ {
			occurrence := domain.KeywordOccurrence{
				PageID:    fmt.Sprintf("page%d", i+1),
				Frequency: 1,
				Positions: []int{1},
			}

			err := repo.Add(keyword, occurrence)

			if err != nil {
				t.Fatalf("Error adding occurrence: %v", err)
			}
		}
	}

	count, err := repo.Count()
	if err != nil {
		t.Fatalf("Error counting occurrences: %v", err)
	}

	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

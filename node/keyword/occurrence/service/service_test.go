package service_test

import (
	"reflect"
	"testing"

	"crawlquery/node/domain"
	"crawlquery/node/keyword/occurrence/repository/mem"
	"crawlquery/node/keyword/occurrence/service"
)

func TestGetKeywordMatches(t *testing.T) {
	repo := mem.NewRepository()
	svc := service.NewKeywordOccurrenceService(repo)

	// Add occurrences to the repository
	keyword := domain.Keyword("example")
	occurrences := []domain.Occurrence{
		{PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
		{PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}
	for _, occ := range occurrences {
		err := repo.Add(keyword, occ)
		if err != nil {
			t.Fatalf("Error adding occurrence: %v", err)
		}
	}

	// Test GetKeywordMatches
	matches, err := svc.GetKeywordMatches([]domain.Keyword{keyword})
	if err != nil {
		t.Fatalf("Error getting keyword matches: %v", err)
	}

	expectedMatches := []domain.KeywordMatch{
		{
			Keyword:     keyword,
			Occurrences: occurrences,
		},
	}

	if !reflect.DeepEqual(matches, expectedMatches) {
		t.Errorf("Expected matches %v, got %v", expectedMatches, matches)
	}
}

func TestUpdateKeywordOccurrences(t *testing.T) {
	repo := mem.NewRepository()
	svc := service.NewKeywordOccurrenceService(repo)

	keyword := domain.Keyword("example")
	occurrence := domain.Occurrence{PageID: "page1", Frequency: 1, Positions: []int{1}}

	err := svc.UpdateKeywordOccurrences("page1", map[domain.Keyword]domain.Occurrence{keyword: occurrence})
	if err != nil {
		t.Fatalf("Error updating keyword occurrences: %v", err)
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

func TestRemovePageOccurrences(t *testing.T) {
	repo := mem.NewRepository()
	svc := service.NewKeywordOccurrenceService(repo)

	keyword := domain.Keyword("example")
	occurrences := []domain.Occurrence{
		{PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
		{PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}
	for _, occ := range occurrences {
		err := repo.Add(keyword, occ)
		if err != nil {
			t.Fatalf("Error adding occurrence: %v", err)
		}
	}

	err := svc.RemovePageOccurrences("page1")
	if err != nil {
		t.Fatalf("Error removing occurrences: %v", err)
	}

	gotOccurrences, err := repo.GetAll(keyword)
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	expectedOccurrences := []domain.Occurrence{
		{PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}

	if !reflect.DeepEqual(gotOccurrences, expectedOccurrences) {
		t.Errorf("Expected occurrences %v, got %v", expectedOccurrences, gotOccurrences)
	}

	err = svc.RemovePageOccurrences("page2")
	if err != nil {
		t.Fatalf("Error removing occurrences: %v", err)
	}

	_, err = repo.GetAll(keyword)
	if err != domain.ErrKeywordNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrKeywordNotFound, err)
	}
}

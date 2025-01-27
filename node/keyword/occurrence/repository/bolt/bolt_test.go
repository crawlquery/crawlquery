package bolt_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"crawlquery/node/domain"
	occRepo "crawlquery/node/keyword/occurrence/repository/bolt"

	"github.com/boltdb/bolt"
)

func setupTestDB(t *testing.T) *bolt.DB {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open BoltDB: %v", err)
	}
	return db
}

func teardownTestDB(db *bolt.DB) {
	db.Close()
	os.Remove("test.db")
}

func TestGet(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo, err := occRepo.NewRepository(db)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	keyword := domain.Keyword("example")
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

	gotOccurrences, err := repo.GetAll(keyword)
	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	if !reflect.DeepEqual(gotOccurrences, occurrences) {
		t.Errorf("Expected occurrences %v, got %v", occurrences, gotOccurrences)
	}

	_, err = repo.GetAll(domain.Keyword("nonexistent"))
	if err != domain.ErrKeywordNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrKeywordNotFound, err)
	}
}

func TestGetForPageID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo, err := occRepo.NewRepository(db)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	keyword := domain.Keyword("example")
	keyword2 := domain.Keyword("example2")
	occurrences := map[domain.Keyword]domain.KeywordOccurrence{
		keyword:  {PageID: "page1", Frequency: 3, Positions: []int{1, 2, 3}},
		keyword2: {PageID: "page2", Frequency: 2, Positions: []int{4, 5}},
	}

	for _, occ := range occurrences {
		err := repo.Add(keyword, occ)
		if err != nil {
			t.Fatalf("Error adding occurrence: %v", err)
		}
	}

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

	empty, err := repo.GetForPageID("nonexistent")

	if err != nil {
		t.Fatalf("Error getting occurrences: %v", err)
	}

	if len(empty) != 0 {
		t.Errorf("Expected 0 occurrences, got %v", len(empty))
	}

}

func TestAddOccurence(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo, err := occRepo.NewRepository(db)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	keyword := domain.Keyword("example")
	occurrence := domain.KeywordOccurrence{PageID: "page1", Frequency: 1, Positions: []int{1}}

	err = repo.Add(keyword, occurrence)
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

func TestRemoveOccurencesForPageID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo, err := occRepo.NewRepository(db)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	keyword := domain.Keyword("example")
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

	err = repo.RemoveForPageID("page1")
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

	err = repo.RemoveForPageID("page2")
	if err != nil {
		t.Fatalf("Error removing occurrences: %v", err)
	}

	_, err = repo.GetAll(keyword)

	if err != domain.ErrKeywordNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrKeywordNotFound, err)
	}
}

func TestCount(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo, err := occRepo.NewRepository(db)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

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

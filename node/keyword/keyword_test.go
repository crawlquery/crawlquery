package keyword_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/keyword"
	"reflect"
	"testing"
)

func TestMakeKeywordOccurrences(t *testing.T) {
	pageID := "page1"
	keywords := []domain.Keyword{
		"example",
		"test",
		"example",
		"sample",
		"example",
		"test",
	}

	expectedOccurrences := map[domain.Keyword]domain.KeywordOccurrence{
		"example": {
			PageID:    pageID,
			Frequency: 3,
			Positions: []int{0, 2, 4},
		},
		"test": {
			PageID:    pageID,
			Frequency: 2,
			Positions: []int{1, 5},
		},
		"sample": {
			PageID:    pageID,
			Frequency: 1,
			Positions: []int{3},
		},
	}

	occurrences, err := keyword.MakeKeywordOccurrences(keywords, pageID)
	if err != nil {
		t.Fatalf("Error making keyword occurrences: %v", err)
	}

	if !reflect.DeepEqual(occurrences, expectedOccurrences) {
		t.Errorf("Expected occurrences %v, got %v", expectedOccurrences, occurrences)
	}
}

package mem

import (
	"crawlquery/api/domain"
	"testing"
	"time"
)

func TestRepo(t *testing.T) {
	// test cases
	tests := []struct {
		// input
		pageID domain.PageID
		rank   float64
		// expected output
		expectedRank float64
	}{
		{"test1", 0.5, 0.5},
		{"test2", 0.3, 0.3},
		{"test3", 0.7, 0.7},
	}

	// setup
	repo := NewRepository()
	for _, test := range tests {
		repo.Update(test.pageID, test.rank, time.Now())
	}

	// test
	for _, test := range tests {
		rank, _ := repo.Get(test.pageID)
		if rank != test.expectedRank {
			t.Errorf("Expected %f, got %f", test.expectedRank, rank)
		}
	}
}

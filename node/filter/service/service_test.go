package service_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/filter/service"
	"testing"
)

func TestPostingsContainsExplicitWords(t *testing.T) {
	t.Run("returns true of contains explicit words", func(t *testing.T) {
		service := service.NewService()

		explicitPostings := map[string]*domain.Posting{
			"fuck": {
				PageID:    "page1",
				Frequency: 2,
				Positions: []int{30, 33},
			},
			"porn": {
				PageID:    "page2",
				Frequency: 2,
				Positions: []int{3},
			},
			"sex": {
				PageID:    "page2",
				Frequency: 4,
				Positions: []int{4},
			},
		}

		if !service.PostingsContainsExplicitWords(explicitPostings) {
			t.Errorf("expected PostingsContainsExplicitWords to return true, got false")
		}
	})

	t.Run("returns false if it does not contain explicit words", func(t *testing.T) {
		service := service.NewService()

		goodWords := map[string]*domain.Posting{
			"hello": {
				PageID:    "page1",
				Frequency: 2,
				Positions: []int{30, 33},
			},
			"person": {
				PageID:    "page2",
				Frequency: 2,
				Positions: []int{3},
			},
			"name": {
				PageID:    "page2",
				Frequency: 4,
				Positions: []int{4},
			},
		}

		if service.PostingsContainsExplicitWords(goodWords) {
			t.Errorf("expected PostingsContainsExplicitWords to return false, got true")
		}
	})
}

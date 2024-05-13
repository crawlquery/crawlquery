package dto_test

import (
	"crawlquery/api/dto"
	"crawlquery/pkg/domain"
	"testing"
)

func TestNewSearchResponse(t *testing.T) {
	t.Run("should create a new search response", func(t *testing.T) {
		// given
		results := []domain.Result{
			{
				PageID: "page1",
				Score:  0.5,
				Page: &domain.Page{
					ID:              "page1",
					URL:             "http://google.com",
					Title:           "Google",
					MetaDescription: "Search the world's information, including webpages, images, videos and more.",
				},
			},
		}

		// when
		res := dto.NewSearchResponse(results)

		// then
		if len(res.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(res.Results))
		}

		if res.Results[0].PageID != results[0].PageID {
			t.Errorf("Expected page ID to be %s, got %s", results[0].PageID, res.Results[0].PageID)
		}

		if res.Results[0].Score != results[0].Score {
			t.Errorf("Expected score to be %f, got %f", results[0].Score, res.Results[0].Score)
		}

		if res.Results[0].Page.ID != results[0].Page.ID {
			t.Errorf("Expected page ID to be %s, got %s", results[0].Page.ID, res.Results[0].Page.ID)
		}

		if res.Results[0].Page.URL != results[0].Page.URL {
			t.Errorf("Expected page URL to be %s, got %s", results[0].Page.URL, res.Results[0].Page.URL)
		}

		if res.Results[0].Page.Title != results[0].Page.Title {
			t.Errorf("Expected page title to be %s, got %s", results[0].Page.Title, res.Results[0].Page.Title)
		}

		if res.Results[0].Page.Description != results[0].Page.MetaDescription {
			t.Errorf("Expected page description to be %s, got %s", results[0].Page.MetaDescription, res.Results[0].Page.Description)
		}
	})
}

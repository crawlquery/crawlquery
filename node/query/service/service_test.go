package service_test

import (
	queryService "crawlquery/node/query/service"
	"testing"
)

func TestQuery(t *testing.T) {
	t.Run("returns results", func(t *testing.T) {
		queryService := queryService.NewService(nil)

		results, err := queryService.Query("SELECT title FROM pages WHERE title LIKE '%example%';")

		if err != nil {
			t.Errorf("Error querying: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		if results[0].ID != "page1" {
			t.Errorf("Expected page ID page1, got %s", results[0].ID)
		}

		if results[0].Title != "Example Page" {
			t.Errorf("Expected title Example Page, got %s", results[0].Title)
		}
	})
}

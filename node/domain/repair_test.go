package domain_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/dto"
	"testing"
	"time"
)

func TestPageDumpFromDTO(t *testing.T) {
	t.Run("can create page dump from DTO", func(t *testing.T) {
		dto := &dto.PageDump{
			PeerID: "peer1",
			PageID: "page1",
			Page: &dto.Page{
				ID:            "page1",
				URL:           "http://google.com",
				Title:         "Google",
				Description:   "Search engine",
				Language:      "English",
				LastIndexedAt: time.Now(),
			},
			KeywordOccurrences: map[string]dto.KeywordOccurrence{
				"keyword1": {
					PageID:    "page1",
					Frequency: 1,
					Positions: []int{1, 2, 3},
				},
			},
		}

		pageDump := domain.PageDumpFromDTO(dto)

		if pageDump.Page.ID != dto.Page.ID {
			t.Errorf("Expected page ID to be %s, got %s", dto.Page.ID, pageDump.Page.ID)
		}

		if pageDump.Page.URL != dto.Page.URL {
			t.Errorf("Expected page URL to be %s, got %s", dto.Page.URL, pageDump.Page.URL)
		}

		if pageDump.Page.Title != dto.Page.Title {
			t.Errorf("Expected page title to be %s, got %s", dto.Page.Title, pageDump.Page.Title)
		}

		if pageDump.Page.Description != dto.Page.Description {
			t.Errorf("Expected page description to be %s, got %s", dto.Page.Description, pageDump.Page.Description)
		}

		if pageDump.Page.Language != dto.Page.Language {
			t.Errorf("Expected page language to be %s, got %s", dto.Page.Language, pageDump.Page.Language)
		}

		if *pageDump.Page.LastIndexedAt != dto.Page.LastIndexedAt {
			t.Errorf("Expected page last indexed at to be %v, got %v", dto.Page.LastIndexedAt, pageDump.Page.LastIndexedAt)
		}

		if len(pageDump.KeywordOccurrences) != 1 {
			t.Errorf("Expected 1 keyword occurrence, got %v", len(pageDump.KeywordOccurrences))
		}

		keywordOccurrence := pageDump.KeywordOccurrences["keyword1"]
		if keywordOccurrence.PageID != "page1" {
			t.Errorf("Expected page ID to be page1, got %s", keywordOccurrence.PageID)
		}

		if keywordOccurrence.Frequency != 1 {
			t.Errorf("Expected frequency to be 1, got %v", keywordOccurrence.Frequency)
		}

		if len(keywordOccurrence.Positions) != 3 {
			t.Errorf("Expected 3 positions, got %v", len(keywordOccurrence.Positions))
		}

		if keywordOccurrence.Positions[0] != 1 {
			t.Errorf("Expected position to be 1, got %v", keywordOccurrence.Positions[0])
		}

		if keywordOccurrence.Positions[1] != 2 {
			t.Errorf("Expected position to be 2, got %v", keywordOccurrence.Positions[1])
		}

		if keywordOccurrence.Positions[2] != 3 {
			t.Errorf("Expected position to be 3, got %v", keywordOccurrence.Positions[2])
		}

	})
}

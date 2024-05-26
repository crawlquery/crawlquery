package domain

import "errors"

var ErrKeywordNotFound = errors.New("keyword not found")

type Keyword string

// Occurrence represents a keyword occurrence in a page.
type Occurrence struct {
	PageID    string `json:"page_id"`
	Frequency int    `json:"frequency"`
	Positions []int  `json:"positions"`
}

type KeywordMatch struct {
	Keyword     Keyword      `json:"keyword"`
	Occurrences []Occurrence `json:"occurrences"`
}

type KeywordOccurrenceService interface {
	GetKeywordMatches(keywords []Keyword) ([]KeywordMatch, error)
	UpdateKeywordOccurrences(pageID string, keywordOccurrences map[Keyword]Occurrence) error
	RemovePageOccurrences(pageID string) error
}

type KeywordOccurrenceRepository interface {
	GetAll(keyword Keyword) ([]Occurrence, error)
	Add(keyword Keyword, occurrence Occurrence) error
	RemoveForPageID(pageID string) error
}

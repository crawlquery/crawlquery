package domain

import "errors"

var ErrKeywordNotFound = errors.New("keyword not found")

type Keyword string

// Occurrence represents a keyword occurrence in a page.
type KeywordOccurrence struct {
	PageID    string `json:"page_id"`
	Frequency int    `json:"frequency"`
	Positions []int  `json:"positions"`
}

type KeywordMatch struct {
	Keyword     Keyword             `json:"keyword"`
	Occurrences []KeywordOccurrence `json:"occurrences"`
}

type KeywordService interface {
	GetKeywordMatches(keywords []Keyword) ([]KeywordMatch, error)
	UpdateOccurrences(pageID string, keywordOccurrences map[Keyword]KeywordOccurrence) error
}

type KeywordOccurrenceRepository interface {
	GetAll(keyword Keyword) ([]KeywordOccurrence, error)
	Add(keyword Keyword, occurrence KeywordOccurrence) error
	RemoveForPageID(pageID string) error
}

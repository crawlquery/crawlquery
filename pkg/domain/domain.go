package domain

import "crawlquery/pkg/index"

type Result struct {
	PageID      string  `json:"id"`
	Url         string  `json:"url"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

type IndexRepository interface {
	Save(idx *index.Index) error
	Load() (*index.Index, error)
}

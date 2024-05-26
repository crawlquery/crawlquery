package domain

import (
	"time"
)

type PageRank struct {
	PageID   string
	PageRank float64
}

type PageRankService interface {
	UpdatePageRanks() error
	GetPageRank(pageID string) (float64, error)
	UpdatePageRanksEvery(duration time.Duration)
}

type PageRankRepository interface {
	Get(pageID string) (float64, error)
	Update(pageID string, rank float64, createdAt time.Time) error
}

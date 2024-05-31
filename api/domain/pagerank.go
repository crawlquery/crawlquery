package domain

import (
	"time"
)

type PageRank struct {
	PageID   PageID
	PageRank float64
}

type PageRankService interface {
	UpdatePageRanks() error
	GetPageRank(pageID PageID) (float64, error)
	UpdatePageRanksEvery(duration time.Duration)
}

type PageRankRepository interface {
	Get(pageID PageID) (float64, error)
	Update(pageID PageID, rank float64, createdAt time.Time) error
}

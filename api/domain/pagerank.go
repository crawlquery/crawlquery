package domain

import "crawlquery/node/domain"

type PageRank struct {
	PageID   string
	PageRank float64
}

type PageRankService interface {
	CalculatePageRank(pageID string) (float64, error)
	ApplyPageRankToResults(results []domain.Result) ([]domain.Result, error)
}

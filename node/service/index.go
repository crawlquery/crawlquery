package service

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/index"
	"crawlquery/pkg/token"
	"sort"
)

type IndexService struct {
	repo  domain.IndexRepository
	index *index.Index
}

func NewIndexService(repo domain.IndexRepository) *IndexService {
	return &IndexService{
		repo: repo,
	}
}

func (service *IndexService) LoadIndex() error {
	idx, err := service.repo.Load()
	if err != nil {
		return err
	}
	service.index = idx
	return nil
}

func (service *IndexService) Search(query string) []domain.Result {
	// Tokenize the query the same way as the index was tokenized
	queryTerms := token.TokenizeTerm(query)
	results := make(map[string]float64) // map[PageID]relevanceScore

	for _, term := range queryTerms {
		if postings, found := service.index.Inverted[term]; found {
			for _, posting := range postings {
				// Simple scoring: count the frequency of each term
				results[posting.PageID] += float64(posting.Frequency)
			}
		}
	}

	// Convert the results map to a slice and sort by relevance score
	var sortedResults []domain.Result
	for docID, score := range results {
		sortedResults = append(sortedResults, domain.Result{PageID: docID, Score: score})
	}
	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].Score > sortedResults[j].Score
	})

	return sortedResults
}

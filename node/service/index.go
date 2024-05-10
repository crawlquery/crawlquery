package service

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/factory"
	"strings"
)

type IndexService struct{}

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (s *IndexService) basicWildcardIndex(term string, results []domain.Result) []domain.Result {
	filteredResults := []domain.Result{}

	// Convert the index term to lowercase for case-insensitive comparison.
	lowerTerm := strings.ToLower(term)

	for _, result := range results {
		// Check if the term is in the title or description, ignoring case.
		if strings.Contains(strings.ToLower(result.Title), lowerTerm) || strings.Contains(strings.ToLower(result.Description), lowerTerm) {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults
}

func (s *IndexService) Index(term string) []domain.Result {
	allPages := factory.ExampleResults()

	if term == "" {
		return allPages
	}

	return s.basicWildcardIndex(term, allPages)
}

package service

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/factory"
	"strings"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) basicWildcardSearch(term string, results []domain.Result) []domain.Result {
	filteredResults := []domain.Result{}

	// Convert the search term to lowercase for case-insensitive comparison.
	lowerTerm := strings.ToLower(term)

	for _, result := range results {
		// Check if the term is in the title or description, ignoring case.
		if strings.Contains(strings.ToLower(result.Title), lowerTerm) || strings.Contains(strings.ToLower(result.Description), lowerTerm) {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults
}

func (s *SearchService) Search(term string) []domain.Result {
	allPages := factory.ExampleResults()

	if term == "" {
		return allPages
	}

	return s.basicWildcardSearch(term, allPages)
}

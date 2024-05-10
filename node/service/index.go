package service

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/index"
	"encoding/gob"
	"os"
	"strings"
)

type IndexService struct {
	index index.Index
}

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (is *IndexService) SaveIndex(filepath string) error {
	// Create a file for writing.
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new gob encoder writing to the file.
	encoder := gob.NewEncoder(file)

	// Encode (serialize) the index.
	if err := encoder.Encode(is.index); err != nil {
		return err
	}

	return nil
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

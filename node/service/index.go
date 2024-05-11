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
	index *index.Index
}

func NewIndexService() *IndexService {
	return &IndexService{
		index: index.NewIndex(),
	}
}

func (is *IndexService) SetIndex(idx *index.Index) {
	is.index = idx
}

func (is *IndexService) GetIndex() *index.Index {
	return is.index
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

// LoadIndex loads an Index from a specified file path using gob encoding.
func (service *IndexService) LoadIndex(filepath string) error {
	// Open the file for reading.
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gob decoder
	decoder := gob.NewDecoder(file)

	// Create an empty Index where the data will be decoded
	var idx index.Index
	if err := decoder.Decode(&idx); err != nil {
		return err
	}
	service.index = &idx

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

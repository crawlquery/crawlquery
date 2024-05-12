package index

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/token"
	"fmt"
	"sort"
	"strings"
)

// Index represents the search index on a node
type Index struct {
	Inverted domain.InvertedIndex
	Forward  domain.ForwardIndex
}

// NewIndex initializes a new Index with prepared structures
func NewIndex() *Index {
	return &Index{
		Inverted: make(domain.InvertedIndex),
		Forward:  make(domain.ForwardIndex),
	}
}

func (idx *Index) SetInverted(inverted domain.InvertedIndex) {
	idx.Inverted = inverted
}

func (idx *Index) SetForward(forward domain.ForwardIndex) {
	idx.Forward = forward
}

func (idx *Index) fuzzySearch(term string) map[string]float64 {
	results := make(map[string]float64)
	for key, postings := range idx.Inverted {
		// Check if the term is a substring of the key
		if strings.Contains(key, term) {
			for _, posting := range postings {
				// Add or increase the score based on frequency
				results[posting.PageID] += float64(posting.Frequency)
			}
		}
	}
	return results
}

func (idx *Index) Search(query string) ([]domain.Result, error) {
	// Tokenize the query the same way as the index was tokenized
	queryTerms := token.TokenizeTerm(query)
	results := make(map[string]float64) // map[PageID]relevanceScore

	for _, term := range queryTerms {
		// Use fuzzy search to find matching terms
		partialResults := idx.fuzzySearch(term)
		for docID, score := range partialResults {
			results[docID] += score
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

	// Add the page metadata to the results
	for i, result := range sortedResults {
		page := idx.Forward[result.PageID]
		sortedResults[i].Page = page
	}

	return sortedResults, nil
}

// AddPage adds a page to both forward and inverted indexes
func (idx *Index) AddPage(doc *domain.Page) {
	tokensWithPositions := token.Tokenize(doc.Content)

	// Update forward index
	idx.Forward[doc.ID] = doc

	// Update inverted index
	for token, positions := range tokensWithPositions {
		posting := domain.Posting{PageID: doc.ID, Frequency: len(positions), Positions: positions}
		idx.Inverted[token] = append(idx.Inverted[token], &posting)
		fmt.Println("adding invertted")
	}

	fmt.Printf("Added %d inverted entries for doc %s\n", len(tokensWithPositions), doc.ID)
}

// Forward returns the forward index
func (idx *Index) GetForward() domain.ForwardIndex {
	return idx.Forward
}

// Inverted returns the inverted index
func (idx *Index) GetInverted() domain.InvertedIndex {
	return idx.Inverted
}
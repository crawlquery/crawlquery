package index

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/token"
	"sort"
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

func (idx *Index) Search(query string) []domain.Result {
	// Tokenize the query the same way as the index was tokenized
	queryTerms := token.TokenizeTerm(query)
	results := make(map[string]float64) // map[PageID]relevanceScore

	for _, term := range queryTerms {
		if postings, found := idx.Inverted[term]; found {
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

// AddDocument adds a document to both forward and inverted indexes
func (idx *Index) AddDocument(doc domain.Document) {
	tokensWithPositions := token.Tokenize(doc.Content)

	// Update forward index
	idx.Forward[doc.ID] = doc

	// Update inverted index
	for token, positions := range tokensWithPositions {
		posting := domain.Posting{PageID: doc.ID, Frequency: len(positions), Positions: positions}
		idx.Inverted[token] = append(idx.Inverted[token], posting)
	}
}

// Forward returns the forward index
func (idx *Index) GetForward() domain.ForwardIndex {
	return idx.Forward
}

// Inverted returns the inverted index
func (idx *Index) GetInverted() domain.InvertedIndex {
	return idx.Inverted
}

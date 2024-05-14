package index

import (
	"crawlquery/node/domain"
	"crawlquery/node/token"
	sharedDomain "crawlquery/pkg/domain"
	"sort"

	"go.uber.org/zap"
)

// Index represents the search index on a node
type Index struct {
	forwardRepo  domain.ForwardIndexRepository
	invertedRepo domain.InvertedIndexRepository
	logger       *zap.SugaredLogger
}

// NewIndex initializes a new Index with prepared structures
func NewIndex(
	forward domain.ForwardIndexRepository,
	inverted domain.InvertedIndexRepository,
	logger *zap.SugaredLogger,
) *Index {
	return &Index{
		forwardRepo:  forward,
		invertedRepo: inverted,
		logger:       logger,
	}
}

func (idx *Index) Search(query string) ([]sharedDomain.Result, error) {
	// Tokenize the query the same way as the index was tokenized
	queryTerms := token.TokenizeTerm(query)
	results := make(map[string]float64) // map[PageID]relevanceScore

	for _, term := range queryTerms {
		// Use fuzzy search to find matching terms
		partialResults := idx.invertedRepo.FuzzySearch(term)
		for docID, score := range partialResults {
			results[docID] += score
		}
	}

	// Convert the results map to a slice and sort by relevance score
	var sortedResults []sharedDomain.Result
	for docID, score := range results {
		sortedResults = append(sortedResults, sharedDomain.Result{PageID: docID, Score: score})
	}

	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].Score > sortedResults[j].Score
	})

	// Add the page metadata to the results
	for i, result := range sortedResults {
		page, err := idx.forwardRepo.Get(result.PageID)

		if err != nil {
			idx.logger.Errorf("Index.Search: Error getting page metadata: %v", err)
			continue
		}
		sortedResults[i].Page = page
	}

	if len(sortedResults) >= 10 {
		sortedResults = sortedResults[:10]
	}

	return sortedResults, nil
}

// AddPage adds a page to both forward and inverted indexes
func (idx *Index) AddPage(doc *sharedDomain.Page) error {
	tokensWithPositions := token.Tokenize(doc.Content)

	// Update forward index
	err := idx.forwardRepo.Save(doc.ID, doc)

	if err != nil {
		idx.logger.Errorf("Index.AddPage: Error saving page metadata: %v", err)

		return err
	}

	// Update inverted index
	for token, positions := range tokensWithPositions {
		posting := domain.Posting{PageID: doc.ID, Frequency: len(positions), Positions: positions}
		idx.invertedRepo.Save(token, &posting)
	}

	return nil
}

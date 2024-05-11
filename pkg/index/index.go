package index

import "crawlquery/pkg/token"

// Document represents a web page with metadata.
type Document struct {
	ID              string
	URL             string
	Title           string
	Content         string // Could be omitted if content isn't stored
	MetaDescription string
}

// Posting lists entry
type Posting struct {
	PageID    string
	Frequency int
	Positions []int // Optional, depending on whether you need positional index
}

// InvertedIndex maps keywords to document lists
type InvertedIndex map[string][]Posting

// ForwardIndex maps document IDs to document metadata and keyword lists
type ForwardIndex map[string]Document

// Index represents the search index on a node
type Index struct {
	Inverted InvertedIndex
	Forward  ForwardIndex
}

// NewIndex initializes a new Index with prepared structures
func NewIndex() *Index {
	return &Index{
		Inverted: make(InvertedIndex),
		Forward:  make(ForwardIndex),
	}
}

// AddDocument adds a document to both forward and inverted indexes
func (idx *Index) AddDocument(doc Document) {
	tokensWithPositions := token.Tokenize(doc.Content)

	// Update forward index
	idx.Forward[doc.ID] = doc

	// Update inverted index
	for token, positions := range tokensWithPositions {
		posting := Posting{PageID: doc.ID, Frequency: len(positions), Positions: positions}
		idx.Inverted[token] = append(idx.Inverted[token], posting)
	}
}

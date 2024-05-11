package domain

type Result struct {
	PageID      string  `json:"id"`
	Url         string  `json:"url"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

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

type Index interface {
	AddDocument(doc Document)
	GetForward() ForwardIndex
	GetInverted() InvertedIndex
	Search(query string) []Result
}

type IndexRepository interface {
	Save(idx Index) error
	Load() (Index, error)
}

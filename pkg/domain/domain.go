package domain

type Result struct {
	PageID string  `json:"id"`
	Score  float64 `json:"score"`
	Page   Page    `json:"page"`
}

// Page represents a web page with metadata.
type Page struct {
	ID              string `json:"id"`
	URL             string `json:"url"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	MetaDescription string `json:"description"`
}

// Posting lists entry
type Posting struct {
	PageID    string
	Frequency int
	Positions []int // Optional, depending on whether you need positional index
}

// InvertedIndex maps keywords to page lists
type InvertedIndex map[string][]Posting

// ForwardIndex maps page IDs to page metadata and keyword lists
type ForwardIndex map[string]Page

type Index interface {
	AddPage(doc Page)
	GetForward() ForwardIndex
	GetInverted() InvertedIndex
	Search(query string) ([]Result, error)
}

type IndexRepository interface {
	Save(idx Index) error
	Load() (Index, error)
}

type IndexService interface {
	Search(query string) ([]Result, error)
}

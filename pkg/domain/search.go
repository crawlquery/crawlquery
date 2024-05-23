package domain

type Result struct {
	PageID string  `json:"id"`
	Score  float64 `json:"score"`
	Page   *Page   `json:"page"`
}

// Page represents a web page with metadata. Note this does not include the keywords.
type Page struct {
	ID              string `json:"id"`
	Hash            string `json:"hash"`
	URL             string `json:"url"`
	Title           string `json:"title"`
	MetaDescription string `json:"meta_description"`
}

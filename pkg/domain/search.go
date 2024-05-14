package domain

type Result struct {
	PageID string  `json:"id"`
	Score  float64 `json:"score"`
	Page   *Page   `json:"page"`
}

// Page represents a web page with metadata.
type Page struct {
	ID               string `json:"id"`
	URL              string `json:"url"`
	Title            string `json:"title"`
	MetaDescription  string `json:"description"`
	Keywords         []string
	KeywordPositions map[string][]int
}

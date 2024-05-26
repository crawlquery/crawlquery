package domain

import "github.com/gin-gonic/gin"

type Result struct {
	PageID            string                       `json:"id"`
	Score             float64                      `json:"score"`
	Page              ResultPage                   `json:"page"`
	KeywordOccurences map[string]KeywordOccurrence `json:"keyword_occurrences"`
	PageRank          float64                      `json:"page_rank"`
}

// Page represents a web page with metadata. Note this does not include the keywords.
type ResultPage struct {
	ID          string `json:"id"`
	Hash        string `json:"hash"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"meta_description"`
}

type SearchService interface {
	Search(query string) ([]Result, error)
}

type SearchHandler interface {
	Search(c *gin.Context)
}

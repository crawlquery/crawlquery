package dto

type CrawlRequest struct {
	PageID string `json:"page_id"`
	URL    string `json:"url"`
}

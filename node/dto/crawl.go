package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type CrawlRequest struct {
	PageID string `json:"page_id"`
	URL    string `json:"url"`
}

type CrawlResponse struct {
	Page *Page `json:"page"`
}

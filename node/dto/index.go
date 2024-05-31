package dto

type IndexRequest struct {
	PageID      string `json:"page_id" binding:"required"`
	URL         string `json:"url" binding:"required"`
	ContentHash string `json:"content_hash" binding:"required"`
}

type IndexResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

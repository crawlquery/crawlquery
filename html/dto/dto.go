package dto

type StorePageRequest struct {
	PageID string `json:"page_id"`
	HTML   []byte `json:"html" binding:"max=10000000"`
}

type StorePageResponse struct {
	Success bool `json:"success"`
}

type GetPageResponse struct {
	HTML []byte `json:"html"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

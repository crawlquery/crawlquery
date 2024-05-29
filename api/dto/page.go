package dto

type CreatePageRequest struct {
	URL string `json:"url"`
}

type CreatePageResponse struct {
	ID string `json:"id"`
}

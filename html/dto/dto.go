package dto

type StorePageRequest struct {
	Hash string `json:"hash"`
	HTML []byte `json:"html" binding:"max=10000000"`
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

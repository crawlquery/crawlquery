package dto

type CreateLinkRequest struct {
	Src string `json:"src" binding:"required,url"`
	Dst string `json:"dst" binding:"required,url"`
}

package handler

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/errorutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	pageService domain.PageService
}

func NewHandler(pageService domain.PageService) *Handler {
	return &Handler{
		pageService: pageService,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := h.pageService.Create(domain.URL(req.URL))

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, dto.CreatePageResponse{
		ID: string(page.ID),
	})
}

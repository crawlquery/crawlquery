package handler

import (
	"crawlquery/node/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service domain.DumpService
}

func NewHandler(service domain.DumpService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Page(c *gin.Context) {
	data, err := h.service.Page()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Data(200, "application/json", data)
}

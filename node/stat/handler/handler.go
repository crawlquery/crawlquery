package handler

import (
	"crawlquery/node/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service domain.StatService
}

func NewHandler(service domain.StatService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Info(c *gin.Context) {
	res, err := h.service.Info()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

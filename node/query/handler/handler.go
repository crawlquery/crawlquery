package handler

import (
	"crawlquery/node/domain"
	"crawlquery/node/dto"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	queryService domain.QueryService
}

func NewHandler(queryService domain.QueryService) *Handler {
	return &Handler{
		queryService: queryService,
	}
}

func (h *Handler) Query(c *gin.Context) {
	var req dto.QueryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var response dto.QueryResponse

	results, err := h.queryService.Query(req.Query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, result := range results {
		response.Pages = append(response.Pages, dto.QueryResultPage{
			ID:    result.ID,
			Title: result.Title,
		})
	}

	c.JSON(200, response)
}

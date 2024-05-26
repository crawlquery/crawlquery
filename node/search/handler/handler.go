package handler

import (
	"crawlquery/node/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SearchHandler struct {
	service domain.SearchService
	logger  *zap.SugaredLogger
}

func NewHandler(service domain.SearchService, logger *zap.SugaredLogger) *SearchHandler {
	return &SearchHandler{
		service: service,
		logger:  logger,
	}
}

func (sh *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(400, gin.H{
			"error": "missing query",
		})
		return
	}

	res, err := sh.service.Search(q)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"results": res,
	})
}

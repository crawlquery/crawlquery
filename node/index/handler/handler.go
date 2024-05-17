package handler

import (
	"crawlquery/node/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IndexHandler struct {
	service domain.IndexService
	logger  *zap.SugaredLogger
}

func NewHandler(service domain.IndexService, logger *zap.SugaredLogger) *IndexHandler {
	return &IndexHandler{
		service: service,
		logger:  logger,
	}
}

func (sh *IndexHandler) Search(c *gin.Context) {
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

func (ih *IndexHandler) Event(c *gin.Context) {
	var event domain.IndexEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		ih.logger.Error(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := ih.service.ApplyIndexEvent(&event); err != nil {
		ih.logger.Error(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "event received",
	})
}

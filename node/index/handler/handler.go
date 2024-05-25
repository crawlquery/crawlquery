package handler

import (
	"crawlquery/node/domain"
	"crawlquery/node/dto"

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

func (ih *IndexHandler) Index(c *gin.Context) {
	pageID := c.Param("pageID")
	if pageID == "" {
		c.JSON(400, &dto.ErrorResponse{
			Error: "missing pageID",
		})
		return
	}

	if err := ih.service.Index(pageID); err != nil {
		ih.logger.Error(err)
		c.JSON(422, &dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(200, &dto.IndexResponse{
		Success: true,
		Message: "indexing complete",
	})
}

func (ih *IndexHandler) GetIndex(c *gin.Context) {
	pageID := c.Param("pageID")
	if pageID == "" {
		c.JSON(400, gin.H{
			"error": "missing pageID",
		})
		return
	}

	index, err := ih.service.GetIndex(pageID)
	if err != nil {
		ih.logger.Error(err)
		c.JSON(404, &dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(200, index)
}

func (ih *IndexHandler) Event(c *gin.Context) {
	var event domain.PageUpdatedEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		ih.logger.Error(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := ih.service.ApplyPageUpdatedEvent(&event); err != nil {
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

func (ih *IndexHandler) Hash(c *gin.Context) {
	pageHash, err := ih.service.Hash()
	if err != nil {
		ih.logger.Error(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"page": pageHash,
	})
}

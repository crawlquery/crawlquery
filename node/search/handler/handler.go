package handler

import (
	"crawlquery/node/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	indexService *service.IndexService
}

func NewSearchHandler(is *service.IndexService) *SearchHandler {
	return &SearchHandler{
		indexService: is,
	}
}

func (sh *SearchHandler) Search(c *gin.Context) {

	res, err := sh.indexService.Search(c.Query("q"))

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

package handler

import (
	"crawlquery/pkg/domain"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	is domain.IndexService
}

func NewSearchHandler(is domain.IndexService) *SearchHandler {
	return &SearchHandler{
		is: is,
	}
}

func (sh *SearchHandler) Search(c *gin.Context) {
	results, err := sh.is.Search(c.Query("q"))

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"results": results,
	})
}

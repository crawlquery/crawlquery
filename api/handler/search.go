package handler

import (
	"crawlquery/pkg/domain"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService domain.SearchService
}

func NewSearchHandler(ss domain.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: ss,
	}
}

func (sh *SearchHandler) Search(c *gin.Context) {
	res, err := sh.searchService.Search(c.Query("q"))
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

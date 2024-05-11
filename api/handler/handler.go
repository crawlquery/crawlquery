package handler

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/shard"
	"net/url"

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

func CrawlHandler(c *gin.Context) {
	if c.Query("url") == "" {
		c.JSON(400, gin.H{
			"error": "url is required",
		})
		return
	}

	// check url is valid
	url, err := url.ParseRequestURI(c.Query("url"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "url is invalid",
		})
		return
	}

	c.JSON(200, gin.H{
		"results": shard.GetShardID(url.String(), 10),
	})
}

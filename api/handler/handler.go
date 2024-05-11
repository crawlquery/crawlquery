package handler

import (
	"crawlquery/api/service"
	"crawlquery/pkg/shard"
	"net/url"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {

	ss := service.NewSearchService()
	c.JSON(200, gin.H{
		"results": ss.Search(c.Query("q")),
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

package handler

import (
	"crawlquery/api/service"
	"crawlquery/pkg/shard"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {

	ss := service.NewSearchService()
	c.JSON(200, gin.H{
		"results": ss.Search(c.Query("q")),
	})
}

func CrawlHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"results": shard.GetShardID(c.Query("url"), 10),
	})
}

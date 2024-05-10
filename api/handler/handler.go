package handler

import (
	"crawlquery/api/service"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {

	ss := service.NewSearchService()
	c.JSON(200, gin.H{
		"results": ss.Search(c.Query("q")),
	})
}

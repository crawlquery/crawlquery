package handler

import (
	"crawlquery/api/service"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {

	is := service.NewIndexService()
	c.JSON(200, gin.H{
		"results": is.Search(c.Query("q")),
	})
}

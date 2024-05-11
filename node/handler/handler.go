package handler

import (
	"crawlquery/node/service"
	"crawlquery/pkg/repository/index/mem"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {

	is := service.NewIndexService(
		mem.NewMemoryRepository(),
	)
	c.JSON(200, gin.H{
		"results": is.Search(c.Query("q")),
	})
}

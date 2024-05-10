package handler

import (
	"crawlquery/pkg/factory"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"results": factory.ExampleResults(),
	})
}

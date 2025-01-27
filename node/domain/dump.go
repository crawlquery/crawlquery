package domain

import "github.com/gin-gonic/gin"

type DumpService interface {
	Page() ([]byte, error)
}

type DumpHandler interface {
	Page(c *gin.Context)
}

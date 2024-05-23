package domain

import "github.com/gin-gonic/gin"

type StatInfo struct {
	TotalPages    int
	TotalKeywords int
	SizeOfIndex   int
}

type StatService interface {
	Info() (*StatInfo, error)
}

type StatHandler interface {
	Info(c *gin.Context)
}

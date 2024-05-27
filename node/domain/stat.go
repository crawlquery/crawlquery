package domain

import "github.com/gin-gonic/gin"

type StatInfo struct {
	TotalPages             int
	TotalIndexedPages      int
	PagesIndexedByLanguage map[string]int
	TotalKeywords          int
	SizeOfPages            int
}

type StatService interface {
	Info() (*StatInfo, error)
}

type StatHandler interface {
	Info(c *gin.Context)
}

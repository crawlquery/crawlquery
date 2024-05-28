package domain

import "github.com/gin-gonic/gin"

type QueryService interface {
	Query(query string) ([]Page, error)
}

type QueryHandler interface {
	Query(c *gin.Context)
}

package domain

import (
	"crawlquery/pkg/domain"

	"github.com/gin-gonic/gin"
)

type SearchService interface {
	Search(term string) ([]domain.Result, error)
}

type SearchHandler interface {
	Search(c *gin.Context)
}

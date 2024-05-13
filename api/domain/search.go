package domain

import (
	"crawlquery/pkg/domain"
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrNoNodesAvailable = errors.New("no nodes available for search")

type SearchService interface {
	Search(term string) ([]domain.Result, error)
}

type SearchHandler interface {
	Search(c *gin.Context)
}

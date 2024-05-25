package domain

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrCrawlFailedToIndexPage = errors.New("failed to index page")

type IndexService interface {
	Search(query string) ([]Result, error)
	Hash() (string, error)
	Index(pageID string) error
	GetIndex(pageID string) (*Page, error)
	ApplyIndexEvent(event *IndexEvent) error
}

type IndexHandler interface {
	Search(c *gin.Context)
	Index(c *gin.Context)
	GetIndex(c *gin.Context)
	Event(c *gin.Context)
	Hash(c *gin.Context)
}

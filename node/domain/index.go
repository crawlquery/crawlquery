package domain

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrCrawlFailedToIndexPage = errors.New("failed to index page")

type IndexService interface {
	Hash() (string, error)
	Index(pageID string, url string, contentHash string) error
	GetIndex(pageID string) (*Page, error)
	ApplyPageUpdatedEvent(event *PageUpdatedEvent) error
}

type IndexHandler interface {
	Index(c *gin.Context)
	GetIndex(c *gin.Context)
	Event(c *gin.Context)
	Hash(c *gin.Context)
}

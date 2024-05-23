package domain

import (
	"crawlquery/pkg/domain"
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrCrawlFailedToIndexPage = errors.New("failed to index page")

// InvertedIndex maps keywords to page lists
type InvertedIndex map[string][]*Posting

// ForwardIndex maps page IDs to page metadata and keyword lists
type ForwardIndex map[string]*domain.Page

type IndexService interface {
	Search(query string) ([]domain.Result, error)
	Hash() (string, error)
	Index(pageID string) error
	ApplyIndexEvent(event *IndexEvent) error
}

type IndexHandler interface {
	Search(c *gin.Context)
	Event(c *gin.Context)
	Hash(c *gin.Context)
}

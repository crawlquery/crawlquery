package domain

import (
	"crawlquery/pkg/domain"

	"github.com/gin-gonic/gin"
)

// Posting lists entry
type Posting struct {
	PageID    string
	Frequency int
	Positions []int // Optional, depending on whether you need positional index
}

// InvertedIndex maps keywords to page lists
type InvertedIndex map[string][]*Posting

// ForwardIndex maps page IDs to page metadata and keyword lists
type ForwardIndex map[string]*domain.Page

type Index interface {
	AddPage(doc *domain.Page)
	GetForward() ForwardIndex
	GetInverted() InvertedIndex
	Search(query string) ([]domain.Result, error)
}

type IndexHandler interface {
	Search(c *gin.Context)
}

type CrawlHandler interface {
	Crawl(c *gin.Context)
}

type CrawlService interface {
	Crawl(url string) error
}

type ForwardIndexRepository interface {
	Get(pageID string) (*domain.Page, error)
	Save(pageID string, page *domain.Page) error
}

type InvertedIndexRepository interface {
	Get(keyword string) ([]*Posting, error)
	Save(token string, posting *Posting) error
	FuzzySearch(token string) map[string]float64
}

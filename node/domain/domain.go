package domain

import (
	"crawlquery/pkg/domain"

	"github.com/gin-gonic/gin"
)

// Posting lists entry
type Posting struct {
	PageID    string `json:"page_id"`
	Frequency int    `json:"frequency"`
	Positions []int  `json:"positions"` // Optional, depending on whether you need positional index
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
	Event(c *gin.Context)
	Hash(c *gin.Context)
}

type CrawlHandler interface {
	Crawl(c *gin.Context)
}

type CrawlService interface {
	Crawl(pageID, url string) error
}

type PageRepository interface {
	Get(pageID string) (*domain.Page, error)
	Save(pageID string, page *domain.Page) error
	Delete(pageID string) error
	GetHashes() (map[string]string, error)
	UpdateHash(pageID, hash string) error
	DeleteHash(pageID string) error
	GetHash(pageID string) (string, error)
}

type PageService interface {
	Get(pageID string) (*domain.Page, error)
	Create(pageID, url string) (*domain.Page, error)
	Update(page *domain.Page) error
	Hash() (string, error)
	JSON() ([]byte, error)
}

type HTMLService interface {
	Get(pageID string) ([]byte, error)
	Save(pageID string, html []byte) error
}

type HTMLRepository interface {
	Get(pageID string) ([]byte, error)
	Save(pageID string, html []byte) error
}

type KeywordRepository interface {
	GetPostings(keyword string) ([]*Posting, error)
	SavePosting(token string, posting *Posting) error
	FuzzySearch(token string) []string
	RemovePostingsByPageID(pageID string) error
	UpdateHash(keyword, hash string) error
	GetHashes() (map[string]string, error)
	GetHash(token string) (string, error)
	GetAll() (map[string][]*Posting, error)
}

type KeywordService interface {
	GetPostings(keyword string) ([]*Posting, error)
	SavePostings(postings map[string]*Posting) error
	FuzzySearch(token string) ([]string, error)
	RemovePostingsByPageID(pageID string) error
	Hash() (string, error)
	JSON() ([]byte, error)
}

type IndexService interface {
	Search(query string) ([]domain.Result, error)
	Hash() (string, string, string, error)
	Index(pageID string) error
	ApplyIndexEvent(event *IndexEvent) error
}

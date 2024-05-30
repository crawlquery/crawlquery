package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrPageNotFound = errors.New("page not found")
var ErrPageAlreadyExists = errors.New("page already exists")
var ErrPageVersionNotFound = errors.New("page version not found")

type PageID string
type URL string
type ContentHash string

type Page struct {
	ID        PageID
	URL       URL
	ShardID   ShardID
	CreatedAt time.Time
}

type PageVersionID string

type PageVersion struct {
	ID          PageVersionID
	PageID      PageID
	ContentHash ContentHash
	CreatedAt   time.Time
}

type PageVersionRepository interface {
	Get(id PageVersionID) (*PageVersion, error)
	Create(v *PageVersion) error
	ListByPageID(pageID PageID) ([]*PageVersion, error)
}

type PageRepository interface {
	Get(id PageID) (*Page, error)
	Create(p *Page) error
}

type PageService interface {
	Get(id PageID) (*Page, error)
	Create(url URL) (*Page, error)
}

type PageHandler interface {
	Create(c *gin.Context)
}

const PageCreatedKey = "page.created"

type PageCreated struct {
	Page *Page
}

func (p PageCreated) Key() EventKey {
	return PageCreatedKey
}

package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrPageNotFound = errors.New("page not found")
var ErrPageAlreadyExists = errors.New("page already exists")

type PageID string
type URL string

type Page struct {
	ID        PageID
	URL       URL
	ShardID   ShardID
	CreatedAt time.Time
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

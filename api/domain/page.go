package domain

import (
	"errors"
	"time"
)

var ErrPageNotFound = errors.New("page not found")
var ErrPageAlreadyExists = errors.New("page already exists")

type PageID string
type URL string
type ShardID uint16

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

package domain

import (
	"errors"
	"time"
)

var ErrPageNotFound = errors.New("page not found")
var ErrPageAlreadyExists = errors.New("page already exists")

type Page struct {
	ID        string
	ShardID   uint
	CreatedAt time.Time
}

type PageRepository interface {
	Get(id string) (*Page, error)
	Create(p *Page) error
}

type PageService interface {
	Get(id string) (*Page, error)
	Create(pageID string, shardID uint) (*Page, error)
}

package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrLinkAlreadyExists = errors.New("link already exists")

type Link struct {
	SrcID     PageID    `json:"src_id"`
	DstID     PageID    `json:"dst_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkRepository interface {
	Create(*Link) error
	GetAll() ([]*Link, error)
}

type LinkService interface {
	Create(srcID PageID, url URL) (*Link, error)
	GetAll() ([]*Link, error)
}

type LinkHandler interface {
	Create(c *gin.Context)
}

const LinkCreatedKey = "link.created"

type LinkCreated struct {
	Link   *Link
	DstURL URL
}

func (l LinkCreated) Key() EventKey {
	return LinkCreatedKey
}

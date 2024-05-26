package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrLinkAlreadyExists = errors.New("link already exists")

type Link struct {
	SrcID     string    `json:"src_id"`
	DstID     string    `json:"dst_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkRepository interface {
	Create(*Link) error
	GetAll() ([]*Link, error)
}

type LinkService interface {
	Create(srcID, dstID string) (*Link, error)
	GetAll() ([]*Link, error)
}

type LinkHandler interface {
	Create(c *gin.Context)
}

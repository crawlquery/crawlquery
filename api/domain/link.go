package domain

import (
	"errors"
	"time"
)

var ErrLinkAlreadyExists = errors.New("link already exists")

type Link struct {
	SrcID     string    `json:"src_id"`
	DstID     string    `json:"dst_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkRepository interface {
	Create(*Link) error
}

type LinkService interface {
	Create(srcID, dstID string) (*Link, error)
}
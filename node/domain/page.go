package domain

import (
	"crawlquery/pkg/domain"
	"errors"
)

var ErrPageNotFound = errors.New("page not found")
var ErrHashNotFound = errors.New("hash not found")

type PageRepository interface {
	Get(pageID string) (*domain.Page, error)
	GetAll() (map[string]*domain.Page, error)
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

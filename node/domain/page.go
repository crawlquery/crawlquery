package domain

import (
	"errors"
	"time"
)

var ErrPageNotFound = errors.New("page not found")
var ErrHashNotFound = errors.New("hash not found")

type Page struct {
	ID            string     `json:"id"`
	Hash          string     `json:"hash"`
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Language      string     `json:"language"`
	LastIndexedAt *time.Time `json:"last_indexed"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type PageRepository interface {
	Get(pageID string) (*Page, error)
	GetAll() (map[string]*Page, error)
	Save(pageID string, page *Page) error
	Delete(pageID string) error
	GetHashes() (map[string]string, error)
	UpdateHash(pageID, hash string) error
	DeleteHash(pageID string) error
	GetHash(pageID string) (string, error)
}

type PageService interface {
	Get(pageID string) (*Page, error)
	GetAll() (map[string]*Page, error)
	Create(pageID, url, hash string) (*Page, error)
	Update(page *Page) error
	UpdateQuietly(page *Page) error
	Hash() (string, error)
	JSON() ([]byte, error)
}

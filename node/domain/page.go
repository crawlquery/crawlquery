package domain

import (
	"errors"
)

var ErrPageNotFound = errors.New("page not found")
var ErrHashNotFound = errors.New("hash not found")

type Page struct {
	ID          string     `json:"id"`
	Hash        string     `json:"hash"`
	URL         string     `json:"url"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Phrases     [][]string `json:"phrases"`
	Language    string     `json:"language"`
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
	Hash() (string, error)
	JSON() ([]byte, error)
}

package domain

import (
	"time"
)

type Result struct {
	PageID string  `json:"id"`
	Score  float64 `json:"score"`
	Page   *Page   `json:"page"`
}

// Page represents a web page with metadata.
type Page struct {
	ID              string `json:"id"`
	URL             string `json:"url"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	MetaDescription string `json:"description"`
}

// Posting lists entry
type Posting struct {
	PageID    string
	Frequency int
	Positions []int // Optional, depending on whether you need positional index
}

// InvertedIndex maps keywords to page lists
type InvertedIndex map[string][]*Posting

// ForwardIndex maps page IDs to page metadata and keyword lists
type ForwardIndex map[string]*Page

type Index interface {
	AddPage(doc *Page)
	GetForward() ForwardIndex
	GetInverted() InvertedIndex
	Search(query string) ([]Result, error)
}

type IndexRepository interface {
	Save(idx Index) error
	Load() (Index, error)
}

type IndexService interface {
	Search(query string) ([]Result, error)
	AddPage(doc *Page) error
}

type CrawlQueueRepository interface {
	Push(j *CrawlJob) error
	Pop() (*CrawlJob, error)
	Save() error
	Load() error
}

type Domain struct {
	Name string
	Lock time.Time
}

type DomainRespository interface {
	Get(name string) (Domain, error)
	Save(d Domain) error
}

type CrawlJob struct {
	URL         string
	RequestedAt time.Time
	LastTriedAt time.Time
	SuccessAt   time.Time
}

type NodeRepository interface {
	Get(id string) (*Node, error)
	GetAll() ([]*Node, error)
	CreateOrUpdate(n *Node) error
	Create(n *Node) error
	Delete(id string) error
}

type NodeService interface {
	Get(id string) (*Node, error)
	GetRandom() (*Node, error)
	RandomizeAll() ([]*Node, error)
	CreateOrUpdate(n *Node) error
	AllByShard() (map[ShardID][]*Node, error)
}

type Node struct {
	ID       string
	Name     string
	ShardID  ShardID
	Hostname string
	Port     string
}

type ShardID int

type SearchService interface {
	Search(term string) ([]Result, error)
}

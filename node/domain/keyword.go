package domain

// Posting lists entry
type Posting struct {
	PageID    string `json:"page_id"`
	Frequency int    `json:"frequency"`
	Positions []int  `json:"positions"` // Optional, depending on whether you need positional index
}

type KeywordRepository interface {
	GetPostings(keyword string) ([]*Posting, error)
	SavePosting(token string, posting *Posting) error
	FuzzySearch(token string) []string
	RemovePostingsByPageID(pageID string) error
	UpdateHash(keyword, hash string) error
	GetHashes() (map[string]string, error)
	GetHash(token string) (string, error)
	GetAll() (map[string][]*Posting, error)
}

type KeywordService interface {
	GetPostings(keyword string) ([]*Posting, error)
	SavePostings(postings map[string]*Posting) error
	FuzzySearch(token string) ([]string, error)
	RemovePostingsByPageID(pageID string) error
	Hash() (string, error)
	JSON() ([]byte, error)
}

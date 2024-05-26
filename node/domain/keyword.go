package domain

// Posting lists entry
type Posting struct {
	PageID    string `json:"page_id"`
	Frequency int    `json:"frequency"`
	Positions []int  `json:"positions"` // Optional, depending on whether you need positional index
}

type KeywordRepository interface {
	GetPages(keyword string) ([]string, error)
	AddPageKeywords(pageID string, keywords []string) error
	RemovePageKeywords(pageID string) error
}

type KeywordService interface {
	UpdatePageKeywords(pageID string, keywords [][]string) error
	GetPageIDsByKeyword(keyword string) ([]string, error)
}

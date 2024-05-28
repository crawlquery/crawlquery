package dto

import (
	"time"
)

type IndexMeta struct {
	PeerID        string    `json:"peer_id"`
	PageID        string    `json:"page_id"`
	LastIndexedAt time.Time `json:"last_indexed_at"`
}

type KeywordOccurrence struct {
	PageID    string `json:"page_id"`
	Frequency int    `json:"frequency"`
	Positions []int  `json:"positions"`
}

type PageDump struct {
	PeerID             string                       `json:"peer_id"`
	PageID             string                       `json:"page_id"`
	Page               Page                         `json:"page"`
	KeywordOccurrences map[string]KeywordOccurrence `json:"keyword_occurences"`
}

type GetIndexMetasResponse struct {
	IndexMetas []IndexMeta `json:"metas"`
}

type GetIndexMetasRequest struct {
	PageIDs []string `json:"page_ids"`
}

type GetPageDumpsRequest struct {
	PageIDs []string `json:"page_ids"`
}

type GetPageDumpsResponse struct {
	PageDumps []PageDump `json:"page_dumps"`
}

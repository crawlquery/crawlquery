package dto

import "time"

type Page struct {
	ID            string    `json:"id"`
	URL           string    `json:"url"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Language      string    `json:"language"`
	Hash          string    `json:"hash"`
	LastIndexedAt time.Time `json:"last_indexed_at"`
}

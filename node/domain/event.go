package domain

import "crawlquery/pkg/domain"

type IndexEvent struct {
	Page     *domain.Page `json:"page"`
	Keywords map[string]*Posting
}

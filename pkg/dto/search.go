package dto

import "crawlquery/pkg/domain"

type NodeSearchResponse struct {
	Results []domain.Result `json:"results"`
}

type SearchResponse NodeSearchResponse

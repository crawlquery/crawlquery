package dto

import "crawlquery/node/domain"

type NodeSearchResponse struct {
	Results []domain.Result `json:"results"`
}

type SearchResponse NodeSearchResponse

package dto

import "crawlquery/pkg/domain"

type SearchResponsePage struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SearchResponseResult struct {
	PageID string             `json:"id"`
	Score  float64            `json:"score"`
	Page   SearchResponsePage `json:"page"`
}

type SearchResponse struct {
	Results []SearchResponseResult `json:"results"`
}

func NewSearchResponse(results []domain.Result) *SearchResponse {
	res := &SearchResponse{}

	for _, r := range results {
		page := SearchResponsePage{
			ID:          r.Page.ID,
			URL:         r.Page.URL,
			Title:       r.Page.Title,
			Description: r.Page.MetaDescription,
		}

		res.Results = append(res.Results, SearchResponseResult{
			PageID: r.PageID,
			Score:  r.Score,
			Page:   page,
		})
	}

	return res
}

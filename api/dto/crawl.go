package dto

import (
	"encoding/json"
	"time"
)

type CreateCrawlJobRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func (r *CreateCrawlJobRequest) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type CrawlJob struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCrawlJobResponse struct {
	CrawlJob CrawlJob `json:"crawl_job"`
}

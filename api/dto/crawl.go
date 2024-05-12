package dto

import (
	"crawlquery/api/domain"
	"encoding/json"
	"time"
)

type CreateCrawlJobRequest struct {
	URL string `json:"url"`
}

func (r *CreateCrawlJobRequest) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type CreateCrawlJobResponse struct {
	CrawlJob struct {
		ID        string    `json:"id"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"crawl_job"`
}

func NewCreateCrawlJobResponse(j *domain.CrawlJob) *CreateCrawlJobResponse {
	res := &CreateCrawlJobResponse{}

	res.CrawlJob.ID = j.ID
	res.CrawlJob.URL = j.URL
	res.CrawlJob.CreatedAt = j.CreatedAt

	return res
}

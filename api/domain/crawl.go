package domain

import "time"

type CrawlJob struct {
	ID        string    `validate:"required,uuid"`
	URL       string    `validate:"required,url"`
	CreatedAt time.Time `validate:"required"`
}

func (j *CrawlJob) Validate() error {
	return validate.Struct(j)
}

type CrawlJobRepository interface {
	Create(*CrawlJob) error
	Get(string) (*CrawlJob, error)
	First() (*CrawlJob, error)
	Delete(string) error
}

type CrawlService interface {
	AddJob(string) error
}

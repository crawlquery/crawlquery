package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrCrawlJobNotFound = errors.New("crawl job not found")

type CrawlJob struct {
	ID            string         `validate:"required,uuid"`
	URL           string         `validate:"required,url"`
	PageID        string         `validate:"required,min=32,max=32"`
	BackoffUntil  sql.NullTime   `validate:""`
	LastCrawledAt sql.NullTime   `validate:""`
	FailedReason  sql.NullString `validate:""`
	CreatedAt     time.Time      `validate:"required"`
}

func (j *CrawlJob) Validate() error {
	return validate.Struct(j)
}

type CrawlJobRepository interface {
	Create(*CrawlJob) error
	Get(string) (*CrawlJob, error)
	Update(*CrawlJob) error
	First() (*CrawlJob, error)
	FirstProcessable() (*CrawlJob, error)
	Delete(string) error
}

type CrawlJobService interface {
	Create(string) (*CrawlJob, error)
	ProcessCrawlJobs()
}

type CrawlJobHandler interface {
	Create(c *gin.Context)
}

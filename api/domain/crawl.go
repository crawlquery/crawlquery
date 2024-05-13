package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrCrawlJobNotFound = errors.New("crawl job not found")

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

type CrawlJobService interface {
	Create(string) (*CrawlJob, error)
	ProcessCrawlJobs()
}

type CrawlJobHandler interface {
	Create(c *gin.Context)
}

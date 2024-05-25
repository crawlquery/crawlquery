package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

// Crawl job errors
var ErrCrawlJobNotFound = errors.New("crawl job not found")

// Lock errors
var ErrInvalidLockKey = errors.New("invalid lock key")
var ErrDomainLocked = errors.New("domain is locked")
var ErrDomainNotLocked = errors.New("domain is not locked")
var ErrCrawlRestrictionNotFound = errors.New("crawl restriction not found")
var ErrCrawlRestrictionAlreadyExists = errors.New("crawl restriction already exists")
var ErrCrawlJobAlreadyExists = errors.New("crawl job already exists")

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
	GetByPageID(string) (*CrawlJob, error)
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

type CrawlRestriction struct {
	Domain string
	Until  sql.NullTime
}

type CrawlRestrictionRepository interface {
	Get(domain string) (*CrawlRestriction, error)
	Set(res *CrawlRestriction) error
}

type CrawlRestrictionService interface {
	GetRestriction(domain string) (bool, *time.Time)
	Restrict(domain string) error
}

package domain

import (
	"errors"
	"time"
)

var ErrCrawlJobNotFound = errors.New("crawl job not found")

type CrawlJobStatus uint8

const (
	CrawlJobStatusPending CrawlJobStatus = iota
	CrawlJobStatusInProgress
	CrawlJobStatusDone
	CrawlJobStatusFailed
)

type CrawlJob struct {
	PageID    PageID
	Status    CrawlJobStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CrawlJobRepository interface {
	Get(pageID PageID) (*CrawlJob, error)
	Save(cj *CrawlJob) error
}

type CrawlLog struct {
	PageID    PageID
	Status    CrawlJobStatus
	Info      string
	CreatedAt time.Time
}

type CrawlLogRepository interface {
	Save(cl *CrawlLog) error
}

type CrawlService interface {
	CreateJob(pageID string) error
}

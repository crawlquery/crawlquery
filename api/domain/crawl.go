package domain

import (
	"errors"
	"time"
)

var ErrCrawlJobNotFound = errors.New("crawl job not found")

type CrawlLogID string
type CrawlStatus uint8

const (
	CrawlStatusPending CrawlStatus = iota
	CrawlStatusInProgress
	CrawlStatusSuccess
	CrawlStatusFailed
)

func (cs CrawlStatus) String() string {
	switch cs {
	case CrawlStatusPending:
		return "pending"
	case CrawlStatusInProgress:
		return "in_progress"
	case CrawlStatusSuccess:
		return "success"
	case CrawlStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type CrawlJob struct {
	PageID    PageID
	Status    CrawlStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CrawlJobRepository interface {
	Get(pageID PageID) (*CrawlJob, error)
	Save(cj *CrawlJob) error
}

type CrawlLog struct {
	ID        CrawlLogID
	PageID    PageID
	Status    CrawlStatus
	Info      string
	CreatedAt time.Time
}

type CrawlLogRepository interface {
	Save(cl *CrawlLog) error
}

type CrawlService interface {
	CreateJob(pageID PageID) error
}

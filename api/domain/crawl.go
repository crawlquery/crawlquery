package domain

import (
	"context"
	"errors"
	"time"
)

var ErrCrawlJobNotFound = errors.New("crawl job not found")
var ErrCrawlQueueEmpty = errors.New("crawl queue is empty")

type CrawlLogID string
type CrawlStatus uint8

const (
	CrawlStatusPending CrawlStatus = iota
	CrawlStatusInProgress
	CrawlStatusCompleted
	CrawlStatusFailed
)

func (cs CrawlStatus) String() string {
	switch cs {
	case CrawlStatusPending:
		return "pending"
	case CrawlStatusInProgress:
		return "in_progress"
	case CrawlStatusCompleted:
		return "completed"
	case CrawlStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type CrawlJob struct {
	PageID    PageID
	URL       URL
	ShardID   ShardID
	Status    CrawlStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CrawlJobRepository interface {
	Get(pageID PageID) (*CrawlJob, error)
	Save(cj *CrawlJob) error
	ListByStatus(limit int, status CrawlStatus) ([]*CrawlJob, error)
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

type CrawlQueue interface {
	Push(job *CrawlJob) error
	Pop() (*CrawlJob, error)
}

type CrawlRateLimiter interface {
	Limit(job *CrawlJob) bool
}

type CrawlService interface {
	CreateJob(page *Page) error
	RunCrawlProcess(ctx context.Context) error
}

type CrawlThrottleService interface {
	CheckAndThrottle(url URL) (bool, error)
}

const CrawlCompletedKey = "crawl.completed"

type CrawlCompleted struct {
	PageID      PageID
	ShardID     ShardID
	ContentHash ContentHash
	Links       []URL
}

func (c CrawlCompleted) Key() EventKey {
	return CrawlCompletedKey
}

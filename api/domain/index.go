package domain

import (
	"database/sql"
	"errors"
	"time"
)

var ErrIndexJobNotFound = errors.New("index job not found")
var ErrIndexJobAlreadyExists = errors.New("index job already exists")

type IndexJob struct {
	ID            string         `json:"id"`
	PageID        string         `json:"page_id"`
	BackoffUntil  sql.NullTime   `json:"backoff_until"`
	LastIndexedAt sql.NullTime   `json:"last_indexed_at"`
	FailedReason  sql.NullString `json:"failed_reason"`
	CreatedAt     time.Time      `json:"created_at"`
}

type IndexJobService interface {
	Get(id string) (*IndexJob, error)
	Next() (*IndexJob, error)
	Create(pageID string) (*IndexJob, error)
	Update(*IndexJob) error
}

type IndexJobRepository interface {
	Get(id string) (*IndexJob, error)
	GetByPageID(pageID string) (*IndexJob, error)
	Next() (*IndexJob, error)
	Create(*IndexJob) (*IndexJob, error)
	Update(*IndexJob) error
}

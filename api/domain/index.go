package domain

import (
	"context"
	"errors"
	"time"
)

var ErrIndexJobNotFound = errors.New("index job not found")
var ErrIndexJobAlreadyExists = errors.New("index job already exists")

type IndexStatus uint8

const (
	IndexStatusPending IndexStatus = iota
	IndexStatusInProgress
	IndexStatusCompleted
	IndexStatusFailed
)

func (is IndexStatus) String() string {
	switch is {
	case IndexStatusPending:
		return "pending"
	case IndexStatusInProgress:
		return "in_progress"
	case IndexStatusCompleted:
		return "completed"
	case IndexStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type IndexJob struct {
	PageID    PageID      `json:"page_id"`
	ShardID   ShardID     `json:"shard_id"`
	Status    IndexStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type IndexService interface {
	CreateJob(pageID PageID, shardID ShardID) error
	RunIndexProcess(ctx context.Context) error
}

type IndexJobRepository interface {
	Get(pageID PageID) (*IndexJob, error)
	Save(*IndexJob) error
	ListByStatus(limit int, status IndexStatus) ([]*IndexJob, error)
}

type IndexLogID string

type IndexLog struct {
	ID        IndexLogID
	PageID    PageID
	Status    IndexStatus
	Info      string
	CreatedAt time.Time
}

type IndexLogRepository interface {
	Save(*IndexLog) error
}

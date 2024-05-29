package domain

import (
	"errors"
	"time"
)

var ErrNoShards = errors.New("no shards")
var ErrShardNotFound = errors.New("shard not found")

type Shard struct {
	ID        uint
	CreatedAt time.Time
}

type ShardRepository interface {
	Count() (int, error)
}

type ShardService interface {
	GetURLShardID(url URL) (ShardID, error)
}

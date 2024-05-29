package domain

import (
	"errors"
	"time"
)

var ErrNoShards = errors.New("no shards")
var ErrShardNotFound = errors.New("shard not found")

type ShardID uint16

type Shard struct {
	ID        ShardID
	CreatedAt time.Time
}

type ShardRepository interface {
	List() ([]*Shard, error)
	Count() (int, error)
}

type ShardService interface {
	List() ([]*Shard, error)
	GetURLShardID(url URL) (ShardID, error)
}

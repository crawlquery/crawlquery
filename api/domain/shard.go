package domain

import (
	"errors"
	"time"
)

var ErrNoShards = errors.New("no shards")

type Shard struct {
	ID        uint
	CreatedAt time.Time
}

type ShardRepository interface {
	Create(*Shard) error
	List() ([]*Shard, error)
	Count() (int, error)
}

type ShardService interface {
	Create(*Shard) error
	List() ([]*Shard, error)
	GetURLShardID(url string) int
}

package domain

import (
	"crawlquery/node/dto"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrInvalidAccountID = errors.New("invalid account id")
var ErrNodeAlreadyExists = errors.New("node already exists")
var ErrNodeNotFound = errors.New("node not found")

type Node struct {
	ID        string    `validate:"required,uuid"`
	Key       string    `validate:"required,uuid"`
	AccountID string    `validate:"required,uuid"`
	Hostname  string    `validate:"required,hostname"`
	Port      uint      `validate:"min=0,max=65535"`
	ShardID   ShardID   `validate:"min=0,max=65535"`
	CreatedAt time.Time `validate:"required"`
}

func (n *Node) Validate() error {
	return validate.Struct(n)
}

type NodeRepository interface {
	Create(*Node) error
	List() ([]*Node, error)
	ListByAccountID(accountID string) ([]*Node, error)
	GetNodeByKey(key string) (*Node, error)
}

type NodeService interface {
	Create(accountID, hostname string, port uint) (*Node, error)
	List() ([]*Node, error)
	RandomizedList() ([]*Node, error)
	RandomizedListGroupByShard() (map[ShardID][]*Node, error)
	ListByAccountID(accountID string) ([]*Node, error)
	ListByShardID(shardID ShardID) ([]*Node, error)
	Randomize(nodes []*Node) []*Node
	SendCrawlJob(node *Node, crawlJob *CrawlJob) (*dto.CrawlResponse, error)
	SendIndexJob(node *Node, indexJob *IndexJob) error
	Auth(key string) (*Node, error)
}

type NodeHandler interface {
	Create(c *gin.Context)
	ListByAccountID(c *gin.Context)
	ListByShardID(c *gin.Context)
	Auth(c *gin.Context)
}

type AllocationService interface {
	AllocateNode(*Node) error
	GetShardWithLeastNodes() (*Shard, error)
}

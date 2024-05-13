package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrInvalidAccountID = errors.New("invalid account id")

type Node struct {
	ID        string    `validate:"required,uuid"`
	AccountID string    `validate:"required,uuid"`
	Hostname  string    `validate:"required,hostname"`
	Port      uint      `validate:"min=0,max=65535"`
	ShardID   uint      `validate:"min=0,max=30000"`
	CreatedAt time.Time `validate:"required"`
}

func (n *Node) Validate() error {
	return validate.Struct(n)
}

type NodeRepository interface {
	Create(*Node) error
	List() ([]*Node, error)
}

type NodeService interface {
	Create(accountID, hostname string, port uint) (*Node, error)
	List() ([]*Node, error)
	RandomizedList() ([]*Node, error)
}

type NodeHandler interface {
	Create(c *gin.Context)
}

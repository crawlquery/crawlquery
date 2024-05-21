package dto

import (
	"crawlquery/api/domain"
	"encoding/json"
	"time"
)

type CreateNodeRequest struct {
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
}

func (r *CreateNodeRequest) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type Node struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	AccountID string    `json:"account_id"`
	Hostname  string    `json:"hostname"`
	Port      uint      `json:"port"`
	ShardID   uint      `json:"shard_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PublicNode struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
	ShardID  uint   `json:"shard_id"`
}

type CreateNodeResponse struct {
	Node *Node `json:"node"`
}

func NewCreateNodeResponse(n *domain.Node) *CreateNodeResponse {
	res := &CreateNodeResponse{
		Node: &Node{
			ID:        n.ID,
			Key:       n.Key,
			AccountID: n.AccountID,
			Hostname:  n.Hostname,
			Port:      n.Port,
			ShardID:   n.ShardID,
			CreatedAt: n.CreatedAt,
		},
	}

	return res
}

type ListNodesResponse struct {
	Nodes []*Node `json:"nodes"`
}

func NewListNodesResponse(nodes []*domain.Node) *ListNodesResponse {
	res := &ListNodesResponse{}

	for _, n := range nodes {
		res.Nodes = append(res.Nodes, &Node{
			ID:        n.ID,
			Key:       n.Key,
			AccountID: n.AccountID,
			Hostname:  n.Hostname,
			Port:      n.Port,
			ShardID:   n.ShardID,
			CreatedAt: n.CreatedAt,
		})
	}

	return res
}

type AuthenticateNodeRequest struct {
	Key string `json:"key" binding:"required"`
}

func NewAuthenticateNodeRequest(key string) *AuthenticateNodeRequest {
	return &AuthenticateNodeRequest{
		Key: key,
	}
}

type AuthenticateNodeResponse struct {
	Node *Node `json:"node"`
}

func NewAuthenticateNodeResponse(n *domain.Node) *AuthenticateNodeResponse {
	res := &AuthenticateNodeResponse{
		Node: &Node{
			ID:        n.ID,
			Key:       n.Key,
			AccountID: n.AccountID,
			Hostname:  n.Hostname,
			Port:      n.Port,
			ShardID:   n.ShardID,
			CreatedAt: n.CreatedAt,
		},
	}

	return res
}

type ListNodesByShardResponse struct {
	Nodes []*PublicNode `json:"nodes"`
}

func NewListNodesByShardResponse(nodes []*domain.Node) *ListNodesByShardResponse {
	res := &ListNodesByShardResponse{}

	for _, n := range nodes {
		res.Nodes = append(res.Nodes, &PublicNode{
			ID:       n.ID,
			Hostname: n.Hostname,
			Port:     n.Port,
			ShardID:  n.ShardID,
		})
	}

	return res
}

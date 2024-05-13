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

type CreateNodeResponse struct {
	Node struct {
		ID        string    `json:"id"`
		AccountID string    `json:"account_id"`
		Hostname  string    `json:"hostname"`
		Port      uint      `json:"port"`
		ShardID   uint      `json:"shard_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"node"`
}

func NewCreateNodeResponse(n *domain.Node) *CreateNodeResponse {
	res := &CreateNodeResponse{}

	res.Node.ID = n.ID
	res.Node.AccountID = n.AccountID
	res.Node.Hostname = n.Hostname
	res.Node.Port = n.Port
	res.Node.ShardID = n.ShardID
	res.Node.CreatedAt = n.CreatedAt

	return res
}

type ListNodesResponse struct {
	Nodes []*struct {
		ID        string    `json:"id"`
		AccountID string    `json:"account_id"`
		Hostname  string    `json:"hostname"`
		Port      uint      `json:"port"`
		ShardID   uint      `json:"shard_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"nodes"`
}

func NewListNodesResponse(nodes []*domain.Node) *ListNodesResponse {
	res := &ListNodesResponse{}

	for _, n := range nodes {
		res.Nodes = append(res.Nodes, &struct {
			ID        string    `json:"id"`
			AccountID string    `json:"account_id"`
			Hostname  string    `json:"hostname"`
			Port      uint      `json:"port"`
			ShardID   uint      `json:"shard_id"`
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:        n.ID,
			AccountID: n.AccountID,
			Hostname:  n.Hostname,
			Port:      n.Port,
			ShardID:   n.ShardID,
			CreatedAt: n.CreatedAt,
		})
	}

	return res
}

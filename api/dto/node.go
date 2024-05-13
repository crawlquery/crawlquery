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

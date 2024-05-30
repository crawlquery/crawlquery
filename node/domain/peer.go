package domain

import (
	"time"
)

type PeerID string

type Peer struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
	ShardID  uint16 `json:"shard_id"`
}

type PageMetadata struct {
	PeerID        string
	PageID        string
	LastIndexedAt time.Time
}

type PeerService interface {
	Self() *Peer
	AddPeer(peer *Peer)
	GetPeers() []*Peer
	GetPeer(id string) (*Peer, error)
	GetIndexMetas(pageIDs []string) ([]IndexMeta, error)
	GetAllIndexMetas() ([]IndexMeta, error)
	GetPageDumpsFromPeer(peer *Peer, pageIDs []PageID) ([]PageDump, error)
	SendPageUpdatedEvent(peer *Peer, event *PageUpdatedEvent) error
	BroadcastPageUpdatedEvent(event *PageUpdatedEvent) error
	SyncPeerList() error
}

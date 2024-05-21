package domain

type Peer struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
	ShardID  uint   `json:"shard_id"`
}

type PeerService interface {
	AddPeer(peer *Peer)
	GetPeers() []*Peer
	GetPeer(id string) *Peer
	SendIndexEvent(peer *Peer, event *IndexEvent) error
	BroadcastIndexEvent(event *IndexEvent) error
	SyncPeerList() error
}

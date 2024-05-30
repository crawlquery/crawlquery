package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/dto"
	"crawlquery/pkg/client/api"
	"crawlquery/pkg/client/node"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	api    *api.Client
	peers  []*domain.Peer
	logger *zap.SugaredLogger
	host   *domain.Peer
	lock   sync.Mutex
}

func NewService(
	api *api.Client,
	host *domain.Peer,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		api:    api,
		host:   host,
		logger: logger,
		lock:   sync.Mutex{},
	}
}

func (s *Service) Self() *domain.Peer {
	return s.host
}

func (s *Service) AddPeer(peer *domain.Peer) {
	if _, err := s.GetPeer(peer.ID); err == nil {
		return
	}

	for _, p := range s.peers {
		if p.Hostname == peer.Hostname && p.Port == peer.Port {
			return
		}
	}

	s.peers = append(s.peers, peer)
}

func (s *Service) GetPeers() []*domain.Peer {
	return s.peers
}

func (s *Service) GetPeer(id string) (*domain.Peer, error) {
	for _, peer := range s.peers {
		if peer.ID == id {
			return peer, nil
		}
	}

	return nil, fmt.Errorf("peer not found")
}

func (s *Service) RemoveAllPeers() {
	s.peers = []*domain.Peer{}
}

func (s *Service) RemovePeer(id string) {
	for i, peer := range s.peers {
		if peer.ID == id {
			s.peers = append(s.peers[:i], s.peers[i+1:]...)
			return
		}
	}
}

func (s *Service) GetPageDumpsFromPeer(peer *domain.Peer, pageIDs []domain.PageID) ([]domain.PageDump, error) {

	client := node.NewClient(
		node.WithHostname(peer.Hostname),
		node.WithPort(peer.Port),
	)

	var strPageIDs []string

	for _, id := range pageIDs {
		strPageIDs = append(strPageIDs, string(id))
	}

	dumps, err := client.GetPageDumps(strPageIDs)

	if err != nil {
		s.logger.Errorf("Error getting page dumps from peer: %v", err)
		return nil, err
	}

	var pageDumps []domain.PageDump

	for _, dump := range dumps {
		pageDumps = append(pageDumps, domain.PageDumpFromDTO(dump))
	}

	return pageDumps, nil
}

func (s *Service) GetIndexMetas(pageIDs []string) ([]domain.IndexMeta, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var wg sync.WaitGroup
	results := make(chan []dto.IndexMeta, len(s.peers))
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent requests

	for _, peer := range s.peers {
		wg.Add(1)
		go func(peer *domain.Peer) {
			defer wg.Done()
			semaphore <- struct{}{}
			metas, err := s.GetIndexMetasFromPeer(peer, pageIDs)
			<-semaphore
			if err != nil {
				results <- nil
			} else {
				results <- metas
			}
		}(peer)
	}

	wg.Wait()
	close(results)

	var allMetas []domain.IndexMeta

	for metas := range results {
		for _, meta := range metas {
			allMetas = append(allMetas, domain.IndexMeta{
				PeerID:        domain.PeerID(meta.PeerID),
				PageID:        domain.PageID(meta.PageID),
				LastIndexedAt: meta.LastIndexedAt,
			})
		}
	}

	return allMetas, nil
}

func (s *Service) GetAllIndexMetas() ([]domain.IndexMeta, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var wg sync.WaitGroup
	results := make(chan []dto.IndexMeta, len(s.peers))
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent requests

	for _, peer := range s.peers {
		wg.Add(1)
		go func(peer *domain.Peer) {
			defer wg.Done()
			semaphore <- struct{}{}
			metas, err := s.GetIndexAllMetasFromPeer(peer)
			<-semaphore
			if err != nil {
				results <- nil
			} else {
				results <- metas
			}
		}(peer)
	}

	wg.Wait()
	close(results)

	var allMetas []domain.IndexMeta

	for metas := range results {
		for _, meta := range metas {
			allMetas = append(allMetas, domain.IndexMeta{
				PeerID:        domain.PeerID(meta.PeerID),
				PageID:        domain.PageID(meta.PageID),
				LastIndexedAt: meta.LastIndexedAt,
			})
		}
	}

	return allMetas, nil
}

func (s *Service) GetIndexAllMetasFromPeer(peer *domain.Peer) ([]dto.IndexMeta, error) {
	client := node.NewClient(
		node.WithHostname(peer.Hostname),
		node.WithPort(peer.Port),
	)

	metas, err := client.GetAllIndexMetas()

	if err != nil {
		s.logger.Errorf("Error getting index metas from peer: %v", err)
		return nil, err
	}

	return metas, nil
}

func (s *Service) GetIndexMetasFromPeer(peer *domain.Peer, pageIDs []string) ([]dto.IndexMeta, error) {
	client := node.NewClient(
		node.WithHostname(peer.Hostname),
		node.WithPort(peer.Port),
	)

	metas, err := client.GetIndexMetas(pageIDs)

	if err != nil {
		s.logger.Errorf("Error getting index metas from peer: %v", err)
		return nil, err
	}

	return metas, nil
}

func (s *Service) SendPageUpdatedEvent(peer *domain.Peer, event *domain.PageUpdatedEvent) error {

	encoded, err := json.Marshal(event)

	if err != nil {
		log.Printf("Error encoding event: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/event", peer.Hostname, peer.Port), bytes.NewBuffer(encoded))

	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		s.logger.Errorf("Error sending event: %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("Error sending event: %s", resp.Status)
		return fmt.Errorf("error sending event: %s", resp.Status)
	}

	s.logger.Infof("Event sent to %s", peer.ID)

	return nil
}

func (s *Service) BroadcastPageUpdatedEvent(event *domain.PageUpdatedEvent) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	var wg sync.WaitGroup
	results := make(chan error, len(s.peers))
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent requests

	for _, peer := range s.peers {
		wg.Add(1)
		go func(peer *domain.Peer) {
			defer wg.Done()
			semaphore <- struct{}{}
			err := s.SendPageUpdatedEvent(peer, event)
			<-semaphore
			results <- err
		}(peer)
	}

	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			log.Printf("Error broadcasting to peer: %v\n", err)
		}
	}

	return nil
}

func (s *Service) SyncPeerList() error {
	nodesInShard, err := s.api.ListNodesByShardID(s.host.ShardID)

	if err != nil {
		s.logger.Errorf("Error listing nodes by shard ID: %v", err)
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.RemoveAllPeers()

	for _, n := range nodesInShard {
		if n.ID == s.host.ID {
			continue
		}
		s.AddPeer(&domain.Peer{
			ID:       n.ID,
			Hostname: n.Hostname,
			Port:     n.Port,
		})
	}

	s.logger.Infof("Synced peer list: %d peers", len(s.peers))
	return nil
}

func (s *Service) SyncPeerListEvery(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		s.SyncPeerList()
	}
}

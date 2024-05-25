package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/pkg/client/api"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	api            *api.Client
	keywordService domain.KeywordService
	pageService    domain.PageService
	peers          []*domain.Peer
	logger         *zap.SugaredLogger
	host           *domain.Peer
	lock           sync.Mutex
}

func NewService(
	api *api.Client,
	pageService domain.PageService,
	host *domain.Peer,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		api:         api,
		pageService: pageService,
		host:        host,
		logger:      logger,
		lock:        sync.Mutex{},
	}
}

func (s *Service) AddPeer(peer *domain.Peer) {
	if s.GetPeer(peer.ID) != nil {
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

func (s *Service) GetPeer(id string) *domain.Peer {
	for _, peer := range s.peers {
		if peer.ID == id {
			return peer
		}
	}

	return nil
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

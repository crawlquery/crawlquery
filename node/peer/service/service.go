package service

import (
	"bytes"
	"crawlquery/node/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type Service struct {
	keywordService domain.KeywordService
	pageService    domain.PageService
	peers          []*domain.Peer
	logger         *zap.SugaredLogger
	host           *domain.Peer
}

func NewService(
	keywordService domain.KeywordService,
	pageService domain.PageService,
	host *domain.Peer,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		keywordService: keywordService,
		pageService:    pageService,
		host:           host,
		logger:         logger,
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

func (s *Service) RemovePeer(id string) {
	for i, peer := range s.peers {
		if peer.ID == id {
			s.peers = append(s.peers[:i], s.peers[i+1:]...)
			return
		}
	}
}

func (s *Service) SendIndexEvent(peer *domain.Peer, event *domain.IndexEvent) error {

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
		return fmt.Errorf("Error sending event: %s", resp.Status)
	}

	s.logger.Infof("Event sent to %s", peer.ID)

	return nil
}

func (s *Service) BroadcastIndexEvent(event *domain.IndexEvent) error {

	var wg sync.WaitGroup
	results := make(chan error, len(s.peers))
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent requests

	for _, peer := range s.peers {
		wg.Add(1)
		go func(peer *domain.Peer) {
			defer wg.Done()
			semaphore <- struct{}{}
			err := s.SendIndexEvent(peer, event)
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

func (s *Service) DiscoverPeers() {

	encoded, err := json.Marshal(s.host)

	if err != nil {
		log.Printf("Error encoding host: %v", err)
		return
	}

	for _, peer := range s.peers {

		endpoint := fmt.Sprintf("http://%s:%d/peers", peer.Hostname, peer.Port)
		fmt.Printf("endpoint: %s\n", endpoint)
		fmt.Printf("encoded: %s\n", encoded)
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(encoded))

		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}

		client := &http.Client{}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			s.logger.Errorf("Error discovering peers: %v", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			s.logger.Errorf("Error discovering peers: %s", resp.Status)
			return
		}

		var peers []*domain.Peer

		err = json.NewDecoder(resp.Body).Decode(&peers)

		if err != nil {
			s.logger.Errorf("Error decoding peers: %v", err)
			return
		}

		for _, p := range peers {
			s.AddPeer(p)
		}
	}
}

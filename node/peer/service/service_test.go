package service_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/peer/service"
	sharedDomain "crawlquery/pkg/domain"
	"crawlquery/pkg/testutil"
	"testing"

	"github.com/h2non/gock"
)

func TestAddPeer(t *testing.T) {
	t.Run("can add peer", func(t *testing.T) {
		service := service.NewService(nil, nil, nil, testutil.NewTestLogger())

		peer := &domain.Peer{
			ID:       "peer1",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}

		service.AddPeer(peer)

		peers := service.GetPeers()

		if len(peers) != 1 {
			t.Fatalf("Expected 1 peer, got %d", len(peers))
		}
	})

	t.Run("can only add the same peer once", func(t *testing.T) {
		service := service.NewService(nil, nil, nil, testutil.NewTestLogger())

		peer := &domain.Peer{
			ID:       "peer1",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}

		service.AddPeer(peer)
		service.AddPeer(peer)

		peers := service.GetPeers()

		if len(peers) != 1 {
			t.Fatalf("Expected 1 peer, got %d", len(peers))
		}
	})
}

func TestGetPeer(t *testing.T) {
	service := service.NewService(nil, nil, nil, testutil.NewTestLogger())

	peer := &domain.Peer{
		ID:       "peer1",
		Hostname: "localhost",
		Port:     8080,
		ShardID:  1,
	}

	service.AddPeer(peer)

	p := service.GetPeer("peer1")

	if p == nil {
		t.Fatalf("Expected to get peer, got nil")
	}

	if p.ID != "peer1" {
		t.Fatalf("Expected peer ID to be peer1, got %s", p.ID)
	}
}

func TestRemovePeer(t *testing.T) {
	service := service.NewService(nil, nil, nil, testutil.NewTestLogger())

	peer := &domain.Peer{
		ID:       "peer1",
		Hostname: "localhost",
		Port:     8080,
		ShardID:  1,
	}

	service.AddPeer(peer)

	service.RemovePeer("peer1")

	peers := service.GetPeers()

	if len(peers) != 0 {
		t.Fatalf("Expected 0 peers, got %d", len(peers))
	}
}

func TestSendIndexEvent(t *testing.T) {
	service := service.NewService(nil, nil, nil, testutil.NewTestLogger())

	peer := &domain.Peer{
		ID:       "peer1",
		Hostname: "localhost",
		Port:     8080,
		ShardID:  1,
	}

	service.AddPeer(peer)

	page := &sharedDomain.Page{
		URL:             "http://example.com",
		ID:              "page1",
		Title:           "Example",
		MetaDescription: "An example page",
	}

	event := &domain.IndexEvent{
		Page: page,
		Keywords: map[string]*domain.Posting{
			"keyword1": {
				PageID:    "page1",
				Frequency: 1,
				Positions: []int{1},
			},
		},
	}

	defer gock.Off()

	gock.New("http://localhost:8080").
		Post("/event").
		JSON(event).
		Reply(200).
		JSON(map[string]interface{}{})

	service.SendIndexEvent(peer, event)

	if !gock.IsDone() {
		t.Fatalf("Expected request to be made")
	}
}

func TestBroadcastIndexEvent(t *testing.T) {
	service := service.NewService(nil, nil, nil, testutil.NewTestLogger())

	peer1 := &domain.Peer{
		ID:       "peer1",
		Hostname: "localhost",
		Port:     8080,
		ShardID:  1,
	}

	peer2 := &domain.Peer{
		ID:       "peer2",
		Hostname: "localhost",
		Port:     8081,
		ShardID:  1,
	}

	service.AddPeer(peer1)
	service.AddPeer(peer2)

	page := &sharedDomain.Page{
		URL:             "http://example.com",
		ID:              "page1",
		Title:           "Example",
		MetaDescription: "An example page",
	}

	event := &domain.IndexEvent{
		Page: page,
		Keywords: map[string]*domain.Posting{
			"keyword1": {
				PageID:    "page1",
				Frequency: 1,
				Positions: []int{1},
			},
		},
	}

	defer gock.Off()

	gock.New("http://localhost:8080").
		Post("/event").
		JSON(event).
		Reply(200).
		JSON(map[string]interface{}{})

	gock.New("http://localhost:8081").
		Post("/event").
		JSON(event).
		Reply(200).
		JSON(map[string]interface{}{})

	service.BroadcastIndexEvent(event)

	if !gock.IsDone() {
		t.Fatalf("Expected requests to be made")
	}
}

func TestDiscoverPeers(t *testing.T) {
	host := &domain.Peer{
		ID:       "peer1",
		Hostname: "localhost",
		Port:     9090,
		ShardID:  1,
	}
	service := service.NewService(nil, nil, host, testutil.NewTestLogger())

	service.AddPeer(&domain.Peer{
		ID:       "peer2",
		Hostname: "localhost",
		Port:     9191,
		ShardID:  1,
	})

	defer gock.Off()

	gock.New("http://localhost:9191").
		Post("/peers").
		JSON(host).
		Reply(200).
		JSON([]domain.Peer{
			{
				ID:       "peer3",
				Hostname: "localhost",
				Port:     9292,
				ShardID:  1,
			},
		})

	service.DiscoverPeers()

	peers := service.GetPeers()

	if len(peers) != 2 {
		t.Fatalf("Expected 2 peers, got %d", len(peers))
	}

	var found bool
	for _, peer := range peers {
		if peer.ID == "peer2" {
			found = true
		}
	}

	if !found {
		t.Fatalf("Expected to find peer2")
	}
}

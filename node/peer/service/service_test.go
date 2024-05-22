package service_test

import (
	"crawlquery/api/dto"
	"crawlquery/node/domain"
	"crawlquery/node/peer/service"
	"crawlquery/pkg/client/api"
	sharedDomain "crawlquery/pkg/domain"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestAddPeer(t *testing.T) {
	t.Run("can add peer", func(t *testing.T) {
		service := service.NewService(nil, nil, nil, nil, testutil.NewTestLogger())

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
		service := service.NewService(nil, nil, nil, nil, testutil.NewTestLogger())

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
	service := service.NewService(nil, nil, nil, nil, testutil.NewTestLogger())

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
	service := service.NewService(nil, nil, nil, nil, testutil.NewTestLogger())

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
	service := service.NewService(nil, nil, nil, nil, testutil.NewTestLogger())

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
	service := service.NewService(nil, nil, nil, nil, testutil.NewTestLogger())

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

func TestRemoveAllPeers(t *testing.T) {
	service := service.NewService(nil, nil, nil, nil, testutil.NewTestLogger())

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

	service.RemoveAllPeers()

	peers := service.GetPeers()

	if len(peers) != 0 {
		t.Fatalf("Expected 0 peers, got %d", len(peers))
	}
}

func TestSyncPeerList(t *testing.T) {

	t.Run("can sync peer list", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://localhost:8080").
			Get("/shards/1/nodes").
			Reply(200).
			JSON(&dto.ListNodesByShardResponse{
				Nodes: []*dto.PublicNode{
					{
						ID:       "peer1",
						Hostname: "localhost",
						Port:     8080,
						ShardID:  1,
					},
					{
						ID:       "peer2",
						Hostname: "localhost",
						Port:     8081,
						ShardID:  1,
					},
				},
			})

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(api, nil, nil, &domain.Peer{
			ID:       "host",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}, testutil.NewTestLogger())

		peer3 := &domain.Peer{
			ID:       "peer3",
			Hostname: "localhost",
			Port:     9191,
			ShardID:  1,
		}

		service.AddPeer(peer3)

		service.SyncPeerList()

		if len(service.GetPeers()) != 2 {
			t.Fatalf("Expected 2 peers, got %d", len(service.GetPeers()))
		}

		if !gock.IsDone() {
			t.Fatalf("Expected request to be made")
		}

		peers := service.GetPeers()

		if peers[0].ID != "peer1" {
			t.Fatalf("Expected peer1, got %s", peers[0].ID)
		}

		if peers[1].ID != "peer2" {
			t.Fatalf("Expected peer2, got %s", peers[1].ID)
		}
	})

	t.Run("can handle error", func(t *testing.T) {
		defer gock.Off()

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(api, nil, nil, &domain.Peer{
			ID:       "host",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}, testutil.NewTestLogger())

		service.SyncPeerList()

		if len(service.GetPeers()) != 0 {
			t.Fatalf("Expected 0 peers, got %d", len(service.GetPeers()))
		}

		if !gock.IsDone() {
			t.Fatalf("Expected request to be made")
		}
	})
}

func TestSyncPeerListEvery(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost:8080").
		Get("/shards/1/nodes").
		Times(10).
		Reply(200).
		JSON(&dto.ListNodesByShardResponse{
			Nodes: []*dto.PublicNode{
				{
					ID:       "peer1",
					Hostname: "localhost",
					Port:     8080,
					ShardID:  1,
				},
			},
		})

	api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

	service := service.NewService(api, nil, nil, &domain.Peer{
		ID:       "host",
		Hostname: "localhost",
		Port:     8080,
		ShardID:  1,
	}, testutil.NewTestLogger())

	go service.SyncPeerListEvery(50 * time.Millisecond)

	time.Sleep(500*time.Millisecond + (30 * time.Nanosecond))

	if len(service.GetPeers()) != 1 {
		t.Fatalf("Expected 1 peer, got %d", len(service.GetPeers()))
	}

	if !gock.IsDone() {
		t.Fatalf("Expected requests to be made")
	}
}

package service_test

import (
	"crawlquery/api/dto"
	"crawlquery/node/domain"
	nodeDto "crawlquery/node/dto"
	"crawlquery/node/peer/service"
	"crawlquery/pkg/client/api"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestAddPeer(t *testing.T) {
	t.Run("can add peer", func(t *testing.T) {
		service := service.NewService(nil, nil, testutil.NewTestLogger())

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
		service := service.NewService(nil, nil, testutil.NewTestLogger())

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
	service := service.NewService(nil, nil, testutil.NewTestLogger())

	peer := &domain.Peer{
		ID:       "peer1",
		Hostname: "localhost",
		Port:     8080,
		ShardID:  1,
	}

	service.AddPeer(peer)

	p, err := service.GetPeer("peer1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if p == nil {
		t.Fatalf("Expected to get peer, got nil")
	}

	if p.ID != "peer1" {
		t.Fatalf("Expected peer ID to be peer1, got %s", p.ID)
	}

	_, err = service.GetPeer("peer2")

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestRemovePeer(t *testing.T) {
	service := service.NewService(nil, nil, testutil.NewTestLogger())

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

func TestSendPageUpdatedEvent(t *testing.T) {
	service := service.NewService(nil, nil, testutil.NewTestLogger())

	peer := &domain.Peer{
		ID:       "peer1",
		Hostname: "localhost",
		Port:     8080,
		ShardID:  1,
	}

	service.AddPeer(peer)

	page := &domain.Page{
		URL:         "http://example.com",
		ID:          "page1",
		Title:       "Example",
		Description: "An example page",
	}

	event := &domain.PageUpdatedEvent{
		Page: page,
	}

	defer gock.Off()

	gock.New("http://localhost:8080").
		Post("/event").
		JSON(event).
		Reply(200).
		JSON(map[string]interface{}{})

	service.SendPageUpdatedEvent(peer, event)

	if !gock.IsDone() {
		t.Fatalf("Expected request to be made")
	}
}

func TestBroadcastPageUpdatedEvent(t *testing.T) {
	service := service.NewService(nil, nil, testutil.NewTestLogger())

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

	page := &domain.Page{
		URL:         "http://example.com",
		ID:          "page1",
		Title:       "Example",
		Description: "An example page",
	}

	event := &domain.PageUpdatedEvent{
		Page: page,
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

	service.BroadcastPageUpdatedEvent(event)

	if !gock.IsDone() {
		t.Fatalf("Expected requests to be made")
	}
}

func TestRemoveAllPeers(t *testing.T) {
	service := service.NewService(nil, nil, testutil.NewTestLogger())

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

		service := service.NewService(api, &domain.Peer{
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

		api := api.NewClient("http://localhost:9202", testutil.NewTestLogger())

		service := service.NewService(api, &domain.Peer{
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

	service := service.NewService(api, &domain.Peer{
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

func TestGetPageDumpsFromPeer(t *testing.T) {
	t.Run("can get page dumps from peer", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://localhost:8080").
			Post("/repair/get-page-dumps").
			Reply(200).
			JSON(&nodeDto.GetPageDumpsResponse{
				PageDumps: []nodeDto.PageDump{
					{
						PeerID: "peer1",
						PageID: "page1",
						Page: nodeDto.Page{
							ID:          "page1",
							URL:         "http://example.com",
							Title:       "Example",
							Description: "An example page",
							Language:    "English",
						},
						KeywordOccurrences: map[string]nodeDto.KeywordOccurrence{
							"keyword1": {
								PageID:    "page1",
								Frequency: 1,
								Positions: []int{1},
							},
						},
					},
				},
			})

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(api, &domain.Peer{
			ID:       "host",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}, testutil.NewTestLogger())

		dumps, err := service.GetPageDumpsFromPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "localhost",
			Port:     8080,
		}, []domain.PageID{"page1"})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(dumps) != 1 {
			t.Fatalf("Expected 1 dump, got %d", len(dumps))
		}

		if dumps[0].PageID != "page1" {
			t.Fatalf("Expected page1, got %s", dumps[0].PageID)
		}

		if dumps[0].PeerID != "peer1" {
			t.Fatalf("Expected peer1, got %s", dumps[0].PeerID)
		}
	})
}

func TestGetIndexMetas(t *testing.T) {
	t.Run("can get index metas", func(t *testing.T) {
		defer gock.Off()

		expectedResponse := &nodeDto.GetIndexMetasResponse{
			IndexMetas: []nodeDto.IndexMeta{
				{
					PeerID:        "peer1",
					PageID:        "page1",
					LastIndexedAt: time.Now(),
				},
			},
		}

		gock.New("http://localhost:8080").
			Post("/repair/get-index-metas").
			Reply(200).
			JSON(expectedResponse)

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(api, &domain.Peer{
			ID:       "host",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}, testutil.NewTestLogger())

		service.AddPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		})

		metas, err := service.GetIndexMetas([]string{"index1"})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(metas) != 1 {
			t.Fatalf("Expected 1 meta, got %d", len(metas))
		}

		if metas[0].PageID != "page1" {
			t.Fatalf("Expected page1, got %s", metas[0].PageID)
		}

		if metas[0].LastIndexedAt.Round(time.Second) != expectedResponse.IndexMetas[0].LastIndexedAt.Round(time.Second) {
			t.Fatalf("Expected last indexed at to be %v, got %v", expectedResponse.IndexMetas[0].LastIndexedAt, metas[0].LastIndexedAt)
		}

		if metas[0].PeerID != "peer1" {
			t.Fatalf("Expected peer1, got %s", metas[0].PeerID)
		}
	})
}

func TestGetAllIndexMetas(t *testing.T) {
	t.Run("can get all index metas", func(t *testing.T) {
		defer gock.Off()

		expectedResponse := &nodeDto.GetIndexMetasResponse{
			IndexMetas: []nodeDto.IndexMeta{
				{
					PeerID:        "peer1",
					PageID:        "page1",
					LastIndexedAt: time.Now(),
				},
			},
		}

		gock.New("http://localhost:8080").
			Get("/repair/get-all-index-metas").
			Reply(200).
			JSON(expectedResponse)

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(api, &domain.Peer{
			ID:       "host",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}, testutil.NewTestLogger())

		service.AddPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		})

		metas, err := service.GetAllIndexMetas()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(metas) != 1 {
			t.Fatalf("Expected 1 meta, got %d", len(metas))
		}

		if metas[0].PageID != "page1" {
			t.Fatalf("Expected page1, got %s", metas[0].PageID)
		}

		if metas[0].LastIndexedAt.Round(time.Second) != expectedResponse.IndexMetas[0].LastIndexedAt.Round(time.Second) {
			t.Fatalf("Expected last indexed at to be %v, got %v", expectedResponse.IndexMetas[0].LastIndexedAt, metas[0].LastIndexedAt)
		}

		if metas[0].PeerID != "peer1" {
			t.Fatalf("Expected peer1, got %s", metas[0].PeerID)
		}
	})
}

func TestGetIndexMetasFromPeer(t *testing.T) {
	t.Run("can get index metas from peer", func(t *testing.T) {
		defer gock.Off()

		expectedResponse := &nodeDto.GetIndexMetasResponse{
			IndexMetas: []nodeDto.IndexMeta{
				{
					PageID:        "page1",
					LastIndexedAt: time.Now(),
				},
			},
		}

		gock.New("http://localhost:8080").
			Post("/repair/get-index-metas").
			Reply(200).
			JSON(expectedResponse)

		api := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		service := service.NewService(api, &domain.Peer{
			ID:       "host",
			Hostname: "localhost",
			Port:     8080,
			ShardID:  1,
		}, testutil.NewTestLogger())

		metas, err := service.GetIndexMetasFromPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "localhost",
			Port:     8080,
		}, []string{"index1"})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(metas) != 1 {
			t.Fatalf("Expected 1 meta, got %d", len(metas))
		}

		if metas[0].PageID != "page1" {
			t.Fatalf("Expected page1, got %s", metas[0].PageID)
		}

		if metas[0].LastIndexedAt.Round(time.Second) != expectedResponse.IndexMetas[0].LastIndexedAt.Round(time.Second) {
			t.Fatalf("Expected last indexed at to be %v, got %v", expectedResponse.IndexMetas[0].LastIndexedAt, metas[0].LastIndexedAt)
		}
	})
}

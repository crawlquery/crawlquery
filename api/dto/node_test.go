package dto_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreateNodeRequestToJSON(t *testing.T) {
	t.Run("should return correct JSON", func(t *testing.T) {
		req := &dto.CreateNodeRequest{
			Hostname: "localhost",
			Port:     8080,
		}

		b, err := req.ToJSON()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := `{"hostname":"localhost","port":8080}`
		if string(b) != expected {
			t.Errorf("expected: %s, got: %s", expected, string(b))
		}
	})
}

func TestNewCreateNodeResponse(t *testing.T) {
	t.Run("should return correct CreateNodeResponse from Node", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUIDString(),
			Key:       util.UUIDString(),
			AccountID: util.UUIDString(),
			Hostname:  "localhost",
			Port:      8080,
			ShardID:   1,
		}

		resp := dto.NewCreateNodeResponse(node)

		if resp.Node.ID != node.ID {
			t.Errorf("expected: %s, got: %s", node.ID, resp.Node.ID)
		}

		if resp.Node.Key != node.Key {
			t.Errorf("expected: %s, got: %s", node.Key, resp.Node.Key)
		}

		if resp.Node.AccountID != node.AccountID {
			t.Errorf("expected: %s, got: %s", node.AccountID, resp.Node.AccountID)
		}

		if resp.Node.Hostname != node.Hostname {
			t.Errorf("expected: %s, got: %s", node.Hostname, resp.Node.Hostname)
		}

		if resp.Node.Port != node.Port {
			t.Errorf("expected: %d, got: %d", node.Port, resp.Node.Port)
		}

		if resp.Node.ShardID != node.ShardID {
			t.Errorf("expected: %d, got: %d", node.ShardID, resp.Node.ShardID)
		}

	})
}

func TestNewListNodesResponse(t *testing.T) {
	t.Run("should return correct ListNodesResponse from Nodes", func(t *testing.T) {
		nodes := []*domain.Node{
			{
				ID:        util.UUIDString(),
				Key:       util.UUIDString(),
				AccountID: util.UUIDString(),
				Hostname:  "localhost",
				Port:      8080,
				ShardID:   1,
			},
			{
				ID:        util.UUIDString(),
				Key:       util.UUIDString(),
				AccountID: util.UUIDString(),
				Hostname:  "localhost",
				Port:      8080,
				ShardID:   1,
			},
		}

		resp := dto.NewListNodesResponse(nodes)

		if len(resp.Nodes) != len(nodes) {
			t.Errorf("expected: %d, got: %d", len(nodes), len(resp.Nodes))
		}

		for i, n := range nodes {
			if resp.Nodes[i].ID != n.ID {
				t.Errorf("expected: %s, got: %s", n.ID, resp.Nodes[i].ID)
			}

			if resp.Nodes[i].Key != n.Key {
				t.Errorf("expected: %s, got: %s", n.Key, resp.Nodes[i].Key)
			}

			if resp.Nodes[i].AccountID != n.AccountID {
				t.Errorf("expected: %s, got: %s", n.AccountID, resp.Nodes[i].AccountID)
			}

			if resp.Nodes[i].Hostname != n.Hostname {
				t.Errorf("expected: %s, got: %s", n.Hostname, resp.Nodes[i].Hostname)
			}

			if resp.Nodes[i].Port != n.Port {
				t.Errorf("expected: %d, got: %d", n.Port, resp.Nodes[i].Port)
			}

			if resp.Nodes[i].ShardID != n.ShardID {
				t.Errorf("expected: %d, got: %d", n.ShardID, resp.Nodes[i].ShardID)
			}
		}
	})
}

func TestNewAuthenticateNodeRequest(t *testing.T) {
	t.Run("should return correct AuthenticateNodeRequest", func(t *testing.T) {
		key := util.UUIDString()
		req := dto.NewAuthenticateNodeRequest(key)

		if req.Key != key {
			t.Errorf("expected: %s, got: %s", key, req.Key)
		}
	})
}

func TestNewAuthenticateNodeResponse(t *testing.T) {
	t.Run("should return correct AuthenticateNodeResponse from Node", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUIDString(),
			Key:       util.UUIDString(),
			AccountID: util.UUIDString(),
			Hostname:  "localhost",
			Port:      8080,
			ShardID:   1,
		}

		resp := dto.NewAuthenticateNodeResponse(node)

		if resp.Node.ID != node.ID {
			t.Errorf("expected: %s, got: %s", node.ID, resp.Node.ID)
		}

		if resp.Node.Key != node.Key {
			t.Errorf("expected: %s, got: %s", node.Key, resp.Node.Key)
		}

		if resp.Node.AccountID != node.AccountID {
			t.Errorf("expected: %s, got: %s", node.AccountID, resp.Node.AccountID)
		}

		if resp.Node.Hostname != node.Hostname {
			t.Errorf("expected: %s, got: %s", node.Hostname, resp.Node.Hostname)
		}

		if resp.Node.Port != node.Port {
			t.Errorf("expected: %d, got: %d", node.Port, resp.Node.Port)
		}

		if resp.Node.ShardID != node.ShardID {
			t.Errorf("expected: %d, got: %d", node.ShardID, resp.Node.ShardID)
		}
	})
}

func TestListNodesByShardResponse(t *testing.T) {
	t.Run("should return correct ListNodesByShardResponse from Nodes", func(t *testing.T) {
		nodes := []*domain.Node{
			{
				ID:        util.UUIDString(),
				Key:       util.UUIDString(),
				AccountID: util.UUIDString(),
				Hostname:  "localhost",
				Port:      8080,
				ShardID:   1,
			},
			{
				ID:        util.UUIDString(),
				Key:       util.UUIDString(),
				AccountID: util.UUIDString(),
				Hostname:  "localhost",
				Port:      8080,
				ShardID:   1,
			},
		}

		resp := dto.NewListNodesByShardResponse(nodes)

		if len(resp.Nodes) != len(nodes) {
			t.Errorf("expected: %d, got: %d", len(nodes), len(resp.Nodes))
		}

		for i, n := range nodes {

			if resp.Nodes[i].ID != n.ID {
				t.Errorf("expected: %s, got: %s", n.ID, resp.Nodes[i].ID)
			}

			if resp.Nodes[i].Hostname != n.Hostname {
				t.Errorf("expected: %s, got: %s", n.Hostname, resp.Nodes[i].Hostname)
			}

			if resp.Nodes[i].Port != n.Port {
				t.Errorf("expected: %d, got: %d", n.Port, resp.Nodes[i].Port)
			}

			if resp.Nodes[i].ShardID != n.ShardID {
				t.Errorf("expected: %d, got: %d", n.ShardID, resp.Nodes[i].ShardID)
			}
		}
	})
}

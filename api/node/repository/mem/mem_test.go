package mem

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create a node", func(t *testing.T) {
		repo := NewRepository()

		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		err := repo.Create(node)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		check, ok := repo.nodes[node.ID]

		if !ok {
			t.Fatalf("Error getting node")
		}

		if check.ID != node.ID {
			t.Errorf("Expected ID to be %s, got %s", node.ID, check.ID)
		}

		if check.AccountID != node.AccountID {
			t.Errorf("Expected AccountID to be %s, got %s", node.AccountID, check.AccountID)
		}

		if check.Hostname != node.Hostname {
			t.Errorf("Expected Hostname to be %s, got %s", node.Hostname, check.Hostname)
		}

		if check.Port != node.Port {
			t.Errorf("Expected Port to be %d, got %d", node.Port, check.Port)
		}

		if check.ShardID != node.ShardID {
			t.Errorf("Expected ShardID to be %d, got %d", node.ShardID, check.ShardID)
		}

		if check.CreatedAt != node.CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, got %v", node.CreatedAt, check.CreatedAt)
		}
	})
}

func TestList(t *testing.T) {
	t.Run("can list nodes", func(t *testing.T) {
		repo := NewRepository()

		node1 := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "testnode1",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		node2 := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "testnode2",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		repo.nodes[node1.ID] = node1
		repo.nodes[node2.ID] = node2

		nodes, err := repo.List()

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(nodes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(nodes))
		}

		if nodes[0].ID != node1.ID {
			t.Errorf("Expected first node to have ID %s, got %s", node1.ID, nodes[0].ID)
		}

		if nodes[1].ID != node2.ID {
			t.Errorf("Expected second node to have ID %s, got %s", node2.ID, nodes[1].ID)
		}
	})
}

func TestListByAccountID(t *testing.T) {
	t.Run("can list nodes by account ID", func(t *testing.T) {
		repo := NewRepository()

		accountID := util.UUID()

		node1 := &domain.Node{
			ID:        util.UUID(),
			AccountID: accountID,
			Hostname:  "testnode1",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		node2 := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "testnode2",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		repo.nodes[node1.ID] = node1
		repo.nodes[node2.ID] = node2

		nodes, err := repo.ListByAccountID(accountID)

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(nodes) != 1 {
			t.Errorf("Expected 1 node, got %d", len(nodes))
		}

		if nodes[0].ID != node1.ID {
			t.Errorf("Expected first node to have ID %s, got %s", node1.ID, nodes[0].ID)
		}
	})
}

package service_test

import (
	"crawlquery/api/service"
	"crawlquery/pkg/domain"
	nodeMemRepo "crawlquery/pkg/repository/node/mem"
	"fmt"
	"testing"
)

func TestNode(t *testing.T) {
	nodeRepo := nodeMemRepo.NewMemoryRepository()
	nodeService := service.NewNodeService(nodeRepo)

	nodeRepo.CreateOrUpdate(&domain.Node{
		ID:       "node1",
		Hostname: "node1.cluster.com",
		Port:     "8080",
	})

	for i := 0; i < 10; i++ {
		nodeRepo.CreateOrUpdate(&domain.Node{
			ID:       fmt.Sprintf("node%d", i),
			Hostname: fmt.Sprintf("node%d.cluster.com", i),
			Port:     "8080",
		})
	}

	node, err := nodeService.Get("node1")
	if err != nil {
		t.Fatalf("Error getting node: %v", err)
	}

	if node.Hostname != "node1.cluster.com" {
		t.Errorf("Expected hostname to be node1.cluster.com, got %v", node.Hostname)
	}

	if node.Port != "8080" {
		t.Errorf("Expected port to be 8080, got %v", node.Port)
	}

	nodes, err := nodeService.RandomizeAll()

	if err != nil {
		t.Fatalf("Error randomizing nodes: %v", err)
	}

	if len(nodes) != 10 {
		t.Errorf("Expected 10 nodes, got %v", len(nodes))
	}

	firstNode := nodes[0]

	timesFirstNodeAppears := 0

	for i := 0; i < 10; i++ {
		nodes, _ := nodeService.RandomizeAll()

		if nodes[0].ID == firstNode.ID {
			timesFirstNodeAppears++
		}
	}

	if timesFirstNodeAppears > 2 {
		t.Errorf("Expected first node to appear only twice, got %v", timesFirstNodeAppears)
	}
}

func TestNodeAllByShard(t *testing.T) {
	nodeRepo := nodeMemRepo.NewMemoryRepository()
	nodeService := service.NewNodeService(nodeRepo)

	for i := 0; i < 10; i++ {
		nodeRepo.CreateOrUpdate(&domain.Node{
			ID:       fmt.Sprintf("node%d", i),
			ShardID:  1,
			Hostname: fmt.Sprintf("node%d.cluster.com", i),
			Port:     "8080",
		})
	}

	for i := 0; i < 5; i++ {
		nodeRepo.CreateOrUpdate(&domain.Node{
			ID:       fmt.Sprintf("node%d", i+10),
			ShardID:  2,
			Hostname: fmt.Sprintf("node%d.cluster.com", i+10),
			Port:     "8080",
		})
	}

	nodes, err := nodeService.AllByShard()

	if err != nil {
		t.Fatalf("Error sharding nodes: %v", err)
	}

	if len(nodes[1]) != 10 {
		t.Errorf("Expected 10 nodes in shard 1, got %v", len(nodes[1]))
	}

	if len(nodes[2]) != 5 {
		t.Errorf("Expected 5 nodes in shard 2, got %v", len(nodes[2]))
	}
}

package disk_test

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/repository/node/disk"
	"os"
	"testing"
)

func TestDisk(t *testing.T) {
	defer os.Remove("/tmp/nodes.gob")
	filepath := "/tmp/nodes.gob"
	repo := disk.NewDiskRepository(filepath)

	err := repo.CreateOrUpdate(&domain.Node{
		ID:       "1",
		Hostname: "node1.cluster.com",
		Port:     "9090",
	})

	if err != nil {
		t.Fatalf("Error pushing to disk repository: %v", err)
	}

	err = repo.CreateOrUpdate(&domain.Node{
		ID:       "2",
		Hostname: "node2.cluster.com",
		Port:     "9090",
	})

	if err != nil {
		t.Fatalf("Error pushing to disk repository: %v", err)
	}

	n, err := repo.Get("1")

	if err != nil {
		t.Fatalf("Error pushing to disk repository: %v", err)
	}

	if n.Hostname != "node1.cluster.com" {
		t.Errorf("Expected URL to be node1.cluster.com, got %v", n.Hostname)
	}

	repoB := disk.NewDiskRepository(filepath)

	repoB.Load()

	n, err = repoB.Get("1")

	if err != nil {
		t.Fatalf("Error fetching from disk repository: %v", err)
	}

	if n.Hostname != "node1.cluster.com" {
		t.Errorf("Expected URL to be node1.cluster.com, got %v", n.Hostname)
	}

	err = repoB.Delete("1")

	if err != nil {
		t.Fatalf("Error pushing to disk repository: %v", err)
	}

	n, err = repoB.Get("1")

	if err == nil {
		t.Fatalf("Expected error getting from disk repository, got nil")
	}

	if n != nil {
		t.Errorf("Expected nil node, got %v", n)
	}

	all, err := repoB.GetAll()

	if err != nil {
		t.Fatalf("Error getting all nodes from disk repository: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected empty list, got %v", all)
	}
}

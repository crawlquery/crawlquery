package disk_test

import (
	"crawlquery/node/html/repository/disk"
	"testing"
)

func TestDisk(t *testing.T) {
	repo, err := disk.NewRepository("/tmp/crawlquery-html")

	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}

	err = repo.Save("test1", []byte("test-data"))

	if err != nil {
		t.Fatalf("Error saving data: %v", err)
	}

	data, err := repo.Get("test1")

	if err != nil {
		t.Fatalf("Error reading data: %v", err)
	}

	if string(data) != "test-data" {
		t.Fatalf("Expected data to be 'test-data', got '%s'", data)
	}
}

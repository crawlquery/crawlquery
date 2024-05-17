package mem_test

import (
	"crawlquery/node/html/repository/mem"
	"testing"
)

func TestMem(t *testing.T) {
	repo := mem.NewRepository()

	err := repo.Save("test1", []byte("test-data"))

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

	_, err = repo.Get("test2")

	if err == nil {
		t.Fatalf("Expected error reading non-existent data")
	}
}

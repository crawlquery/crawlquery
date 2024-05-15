package service_test

import (
	"crawlquery/node/html/repository/mem"
	"crawlquery/node/html/service"
	"testing"
)

func Test(t *testing.T) {
	repo := mem.NewRepository()

	service := service.NewService(repo)

	err := service.Save("test1", []byte("test-data"))

	if err != nil {
		t.Fatalf("Error saving data: %v", err)
	}

	data, err := service.Get("test1")

	if err != nil {
		t.Fatalf("Error reading data: %v", err)
	}

	if string(data) != "test-data" {
		t.Fatalf("Expected data to be 'test-data', got '%s'", data)
	}

	_, err = service.Get("test2")

	if err == nil {
		t.Fatalf("Expected error reading non-existent data")
	}
}

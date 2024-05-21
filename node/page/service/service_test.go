package service_test

import (
	pageRepo "crawlquery/node/page/repository/mem"
	"crawlquery/node/page/service"
	"testing"
)

func TestCreate(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	page, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	if page.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", page.ID)
	}

	if page.URL != "http://example.com" {
		t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", page.URL)
	}

	check, err := pageRepo.Get("1")

	if err != nil {
		t.Fatalf("Error getting page: %v", err)
	}

	if check.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", check.ID)
	}

	if check.URL != "http://example.com" {
		t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", check.URL)
	}
}

func TestGet(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	_, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	check, err := service.Get("1")

	if err != nil {
		t.Fatalf("Error getting page: %v", err)
	}

	if check.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", check.ID)
	}

	if check.URL != "http://example.com" {
		t.Fatalf("Expected page URL to be 'http://example.com', got '%s'", check.URL)
	}
}

func TestUpdate(t *testing.T) {
	pageRepo := pageRepo.NewRepository()
	service := service.NewService(pageRepo)

	page, err := service.Create("1", "http://example.com")

	if err != nil {
		t.Fatalf("Error saving page: %v", err)
	}

	page.URL = "http://example2.com"

	err = service.Update(page)

	if err != nil {
		t.Fatalf("Error updating page: %v", err)
	}

	check, err := service.Get("1")

	if err != nil {
		t.Fatalf("Error getting page: %v", err)
	}

	if check.ID != "1" {
		t.Fatalf("Expected page ID to be '1', got '%s'", check.ID)
	}

	if check.URL != "http://example2.com" {
		t.Fatalf("Expected page URL to be 'http://example2.com', got '%s'", check.URL)
	}
}
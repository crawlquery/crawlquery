package service_test

import (
	"crawlquery/html/dto"
	"crawlquery/node/html/repository/mem"
	"crawlquery/node/html/service"
	"crawlquery/pkg/client/html"

	backupService "crawlquery/node/html/backup/service"
	"testing"

	"github.com/h2non/gock"
)

func TestGet(t *testing.T) {
	t.Run("returns page", func(t *testing.T) {

		repo := mem.NewRepository()

		backupService := backupService.NewService(html.NewClient("http://storage:8080"))

		service := service.NewService(repo, backupService)

		err := repo.Save("test1", []byte("test-data"))

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
	})

	t.Run("returns backup data", func(t *testing.T) {
		defer gock.Off()

		content := []byte("<html><body><h1>Hello, World!</h1></body></html>")

		resp := dto.GetPageResponse{
			HTML: content,
		}

		gock.New("http://storage:8080").
			Get("/pages/page1").
			Reply(200).
			JSON(resp)

		client := html.NewClient("http://storage:8080")

		service := backupService.NewService(client)

		page, err := service.Get("page1")

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if string(page) != string(content) {
			t.Fatalf("Test failed: Expected body %v, got %v", string(content), string(page))
		}
	})

	t.Run("returns error if not found", func(t *testing.T) {
		repo := mem.NewRepository()

		backupService := backupService.NewService(html.NewClient("http://storage:8080"))

		service := service.NewService(repo, backupService)

		_, err := service.Get("not-found")

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestSave(t *testing.T) {
	t.Run("stores page", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://storage:8080").
			Post("/pages").
			Reply(201)

		repo := mem.NewRepository()

		backupService := backupService.NewService(html.NewClient("http://storage:8080"))

		service := service.NewService(repo, backupService)

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
	})
}

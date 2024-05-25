package service_test

import (
	"crawlquery/html/dto"
	"crawlquery/pkg/client/html"

	backupService "crawlquery/node/html/backup/service"
	"testing"

	"github.com/h2non/gock"
)

func TestGet(t *testing.T) {
	t.Run("returns page", func(t *testing.T) {
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
}

func TestSave(t *testing.T) {
	t.Run("stores page", func(t *testing.T) {
		defer gock.Off()

		expectedReq := dto.StorePageRequest{
			PageID: "page1",
			HTML:   []byte("<html><body><h1>Hello, World!</h1></body></html>"),
		}

		gock.New("http://storage:8080").
			Post("/pages").
			JSON(expectedReq).
			Reply(200)

		client := html.NewClient("http://storage:8080")

		service := backupService.NewService(client)

		err := service.Save("page1", []byte("<html><body><h1>Hello, World!</h1></body></html>"))

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}
	})
}

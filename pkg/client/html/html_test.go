package html_test

import (
	"crawlquery/html/dto"
	"crawlquery/pkg/client/html"
	"testing"

	"github.com/h2non/gock"
)

func TestGetPage(t *testing.T) {
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

		page, err := client.GetPage("page1")

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if string(page) != string(content) {
			t.Fatalf("Test failed: Expected body %v, got %v", string(content), string(page))
		}
	})
}

func TestStorePage(t *testing.T) {
	t.Run("stores page", func(t *testing.T) {
		defer gock.Off()

		expectedReq := dto.StorePageRequest{
			Hash: "page1",
			HTML: []byte("<html><body><h1>Hello, World!</h1></body></html>"),
		}

		gock.New("http://storage:8080").
			Post("/pages").
			JSON(expectedReq).
			Reply(200)

		client := html.NewClient("http://storage:8080")

		err := client.StorePage("page1", []byte("<html><body><h1>Hello, World!</h1></body></html>"))

		if err != nil {
			t.Fatalf("Test failed: %v", err)
		}

		if !gock.IsDone() {
			t.Fatalf("Test failed: Expected request not made")
		}
	})
}

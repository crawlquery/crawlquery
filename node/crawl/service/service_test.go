package service_test

import (
	"crawlquery/node/crawl/service"
	"crawlquery/node/html/repository/disk"
	"crawlquery/pkg/testutil"
	"os"
	"testing"

	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	t.Run("can crawl a page", func(t *testing.T) {

		path := "/tmp/crawlquery-node-crawl-test"
		defer os.RemoveAll(path)
		htmlRepo, err := disk.NewRepository(path)
		if err != nil {
			t.Fatalf("Error creating repository: %v", err)
		}

		service := service.NewService(htmlRepo, testutil.NewTestLogger())

		defer gock.Off()

		expectedData := "<html><head><title>Example</title></head><body><h1>Hello, World!</h1></body></html>"

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString(expectedData)

		err = service.Crawl("test1", "http://example.com")

		if err != nil {
			t.Errorf("Error crawling page: %v", err)
		}

		if err != nil {
			t.Fatalf("Error creating repository: %v", err)
		}

		data, err := htmlRepo.Read("test1")

		if err != nil {
			t.Fatalf("Error reading data: %v", err)
		}

		if string(data) != expectedData {
			t.Fatalf("Expected data to be '%s', got '%s'", expectedData, data)
		}
	})

	t.Run("handles error saving page", func(t *testing.T) {
		path := "/tmp/crawlquery-node-crawl-test"
		defer os.RemoveAll(path)
		htmlRepo, err := disk.NewRepository(path)
		if err != nil {
			t.Fatalf("Error creating repository: %v", err)
		}

		service := service.NewService(htmlRepo, testutil.NewTestLogger())

		defer gock.Off()

		gock.New("http://example.com").
			Get("/").
			Reply(200).
			BodyString("")

		err = service.Crawl("test1", "http://example2.com")

		if err == nil {
			t.Fatalf("Expected error crawling page")
		}

		data, err := htmlRepo.Read("test1")

		if err == nil {
			t.Fatalf("Expected error reading data")
		}

		if data != nil {
			t.Fatalf("Expected data to be nil")
		}
	})
}

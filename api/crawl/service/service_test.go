package service_test

import (
	"crawlquery/api/service"
	"crawlquery/pkg/domain"
	crawlJobMemRepo "crawlquery/pkg/repository/job/mem"
	nodeMemRepo "crawlquery/pkg/repository/node/mem"
	"testing"

	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	nodeRepo := nodeMemRepo.NewMemoryRepository()
	nodeService := service.NewNodeService(nodeRepo)
	crawlJobRepo := crawlJobMemRepo.NewMemoryRepository()

	nodeRepo.CreateOrUpdate(&domain.Node{
		ID:       "node1",
		Hostname: "node1.cluster.com",
		Port:     "8080",
	})

	crawlJobRepo.Push(&domain.CrawlJob{
		URL: "http://google.com",
	})

	svc := service.NewCrawlService(crawlJobRepo, nodeService)

	defer gock.Off()

	gock.New("http://node1.cluster.com:8080").
		MatchHeader("Content-Type", "application/json").
		Post("/crawl").
		JSON(map[string]string{"url": "http://google.com"}).
		Reply(200).
		JSON(map[string]bool{"success": true})

	err := svc.Crawl()

	if err != nil {
		t.Fatalf("Error crawling: %v", err)
	}
}

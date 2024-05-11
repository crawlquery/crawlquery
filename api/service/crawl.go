package service

import (
	"bytes"
	"crawlquery/pkg/domain"
	"fmt"
	"net/http"
)

type CrawlService struct {
	queueRepository domain.CrawlQueueRepository
	nodeService     domain.NodeService
}

func NewCrawlService(
	queueRepository domain.CrawlQueueRepository,
	nodeService domain.NodeService,
) *CrawlService {
	return &CrawlService{
		queueRepository: queueRepository,
		nodeService:     nodeService,
	}
}

func (service *CrawlService) Crawl() error {
	err := service.queueRepository.Load()
	if err != nil {
		return err
	}

	job, err := service.queueRepository.Pop()

	if err != nil {
		return err
	}
	nodes, err := service.nodeService.RandomizeAll()

	for _, node := range nodes {
		if err != nil {
			return err
		}

		crawlEndpoint := fmt.Sprintf("http://%s:%d/crawl", node.Hostname, node.Port)

		res, err := http.Post(crawlEndpoint, "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(
			`{"url": "%s"}`, job.URL,
		))))

		if err != nil {
			return err
		}

		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			break
		}
	}

	return nil
}

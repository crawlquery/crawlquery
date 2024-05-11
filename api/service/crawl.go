package service

import (
	"bytes"
	"crawlquery/pkg/domain"
	"fmt"
	"net/http"
	"time"
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

func (service *CrawlService) Queue(url string) error {
	return service.queueRepository.Push(
		&domain.CrawlJob{
			URL:         url,
			RequestedAt: time.Now(),
		},
	)
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
	fmt.Printf("nodes: %v\n", nodes)
	var finished bool
	for _, node := range nodes {
		if err != nil {
			return err
		}

		crawlEndpoint := fmt.Sprintf("http://%s:%s/crawl", node.Hostname, node.Port)

		res, err := http.Post(crawlEndpoint, "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(
			`{"url": "%s"}`, job.URL,
		))))

		if err != nil {
			return err
		}

		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			finished = true
			break
		}
	}

	if !finished {
		job.LastTriedAt = time.Now()

		err = service.queueRepository.Push(job)

		if err != nil {
			return err
		}

		return fmt.Errorf("could not crawl %s", job.URL)
	}

	return nil
}

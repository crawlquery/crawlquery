package mem_test

import (
	"crawlquery/pkg/util"
	"testing"

	"crawlquery/api/crawl/queue/mem"
	"crawlquery/api/domain"
)

func TestMemQueue(t *testing.T) {
	// test cases
	tests := []struct {
		// input
		crawlJob *domain.CrawlJob
		// expected output
		expectedPageID domain.PageID
	}{
		{
			&domain.CrawlJob{
				PageID: util.PageID("http://example.org"),
			},
			util.PageID("http://example.org"),
		},
		{
			&domain.CrawlJob{
				PageID: util.PageID("http://example.net"),
			},
			util.PageID("http://example.net"),
		},
		{
			&domain.CrawlJob{
				PageID: util.PageID("http://example.com"),
			},
			util.PageID("http://example.com"),
		},
	}

	// setup
	repo := mem.NewQueue()
	for _, test := range tests {
		repo.Push(test.crawlJob)
	}

	// test
	for _, test := range tests {
		job, _ := repo.Pop()
		if job.PageID != test.expectedPageID {
			t.Errorf("Expected %v, got %v", test.expectedPageID, job.PageID)
		}
	}

	// test empty queue
	job, err := repo.Pop()

	if err != domain.ErrCrawlQueueEmpty {
		t.Errorf("Expected error, got %v", err)
	}

	if job != nil {
		t.Errorf("Expected job to be nil, got %v", job)
	}
}

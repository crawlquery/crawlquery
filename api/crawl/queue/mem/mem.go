package mem

import (
	"crawlquery/api/domain"
	"sync"
)

type Queue struct {
	jobs  []*domain.CrawlJob
	mutex *sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		jobs:  make([]*domain.CrawlJob, 0),
		mutex: &sync.Mutex{},
	}
}

func (r *Queue) Push(job *domain.CrawlJob) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.jobs = append(r.jobs, job)
	return nil
}

func (r *Queue) Pop() (*domain.CrawlJob, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.jobs) == 0 {
		return nil, domain.ErrCrawlQueueEmpty
	}

	job := r.jobs[0]
	r.jobs = r.jobs[1:]

	return job, nil
}

package mem

import (
	"crawlquery/api/domain"
	"sync"
)

type Repository struct {
	jobs  []*domain.CrawlJob
	mutex *sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		jobs:  make([]*domain.CrawlJob, 0),
		mutex: &sync.Mutex{},
	}
}

func (r *Repository) Push(job *domain.CrawlJob) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.jobs = append(r.jobs, job)
	return nil
}

func (r *Repository) Pop() (*domain.CrawlJob, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.jobs) == 0 {
		return nil, nil
	}

	job := r.jobs[0]
	r.jobs = r.jobs[1:]

	return job, nil
}

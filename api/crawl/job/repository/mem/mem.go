package mem

import (
	"crawlquery/api/domain"
	"sync"
)

type Repository struct {
	jobs  map[domain.PageID]*domain.CrawlJob
	mutex *sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		jobs:  make(map[domain.PageID]*domain.CrawlJob),
		mutex: &sync.Mutex{},
	}
}

func (r *Repository) Get(pageID domain.PageID) (*domain.CrawlJob, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	job, ok := r.jobs[pageID]
	if !ok {
		return nil, domain.ErrCrawlJobNotFound
	}
	return job, nil
}

func (r *Repository) Save(cj *domain.CrawlJob) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.jobs[cj.PageID] = cj
	return nil
}

func (r *Repository) ListByStatus(limit int, status domain.CrawlStatus) ([]*domain.CrawlJob, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var jobs []*domain.CrawlJob
	for _, job := range r.jobs {
		if job.Status == status {
			jobs = append(jobs, job)
		}
	}
	if len(jobs) > limit {
		jobs = jobs[:limit]
	}
	return jobs, nil
}

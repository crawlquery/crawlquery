package mem

import (
	"crawlquery/api/domain"
	"time"
)

type Repository struct {
	jobs       map[string]*domain.CrawlJob
	forceError error
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make(map[string]*domain.CrawlJob),
	}
}

func (r *Repository) ForceError(err error) {
	r.forceError = err
}

func (r *Repository) Create(j *domain.CrawlJob) error {
	if r.forceError != nil {
		return r.forceError
	}
	r.jobs[j.ID] = j
	return nil
}

func (r *Repository) Update(j *domain.CrawlJob) error {
	if r.forceError != nil {
		return r.forceError
	}

	r.jobs[j.ID] = j
	return nil
}

func (r *Repository) Get(id string) (*domain.CrawlJob, error) {
	job, ok := r.jobs[id]
	if !ok {
		return nil, domain.ErrCrawlJobNotFound
	}

	return job, nil
}

func (r *Repository) GetByPageID(pageID string) (*domain.CrawlJob, error) {
	for _, job := range r.jobs {
		if job.PageID == pageID {
			return job, nil
		}
	}

	return nil, domain.ErrCrawlJobNotFound
}

func (r *Repository) First() (*domain.CrawlJob, error) {

	var earliest *domain.CrawlJob
	for _, job := range r.jobs {

		if earliest == nil || job.CreatedAt.Before(earliest.CreatedAt) {
			earliest = job
		}
	}

	return earliest, nil
}

func (r *Repository) FirstProcessable() (*domain.CrawlJob, error) {
	var earliest *domain.CrawlJob
	for _, job := range r.jobs {
		if job.BackoffUntil.Valid && job.BackoffUntil.Time.After(job.CreatedAt) {
			continue
		}

		if job.LastCrawledAt.Valid && !job.LastCrawledAt.Time.Add(time.Hour*24*31).Before(time.Now()) {
			continue
		}

		if earliest == nil || job.CreatedAt.Before(earliest.CreatedAt) {
			earliest = job
		}
	}

	return earliest, nil
}

func (r *Repository) Delete(id string) error {
	delete(r.jobs, id)
	return nil
}

package mem

import (
	"crawlquery/api/domain"
	"time"
)

type Repository struct {
	jobs map[string]*domain.IndexJob
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make(map[string]*domain.IndexJob),
	}
}

func (r *Repository) Get(id string) (*domain.IndexJob, error) {
	job, ok := r.jobs[id]

	if !ok {
		return nil, domain.ErrIndexJobNotFound
	}

	return job, nil
}

func (r *Repository) GetByPageID(pageID string) (*domain.IndexJob, error) {
	for _, job := range r.jobs {
		if job.PageID == pageID {
			return job, nil
		}
	}

	return nil, domain.ErrIndexJobNotFound
}

func (r *Repository) Next() (*domain.IndexJob, error) {
	for _, job := range r.jobs {

		if job.BackoffUntil.Valid && !job.BackoffUntil.Time.After(time.Now()) {
			continue
		}

		if job.LastIndexedAt.Valid {
			continue
		}

		return job, nil
	}

	return nil, domain.ErrIndexJobNotFound
}

func (r *Repository) Create(job *domain.IndexJob) (*domain.IndexJob, error) {

	r.jobs[job.ID] = job

	return job, nil
}

func (r *Repository) Update(job *domain.IndexJob) error {
	r.jobs[job.ID] = job

	return nil
}

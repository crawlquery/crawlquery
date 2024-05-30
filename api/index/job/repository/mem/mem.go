package mem

import (
	"crawlquery/api/domain"
)

type Repository struct {
	jobs map[domain.PageID]*domain.IndexJob
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make(map[domain.PageID]*domain.IndexJob),
	}
}

func (r *Repository) Get(pageID domain.PageID) (*domain.IndexJob, error) {
	job, ok := r.jobs[pageID]

	if !ok {
		return nil, domain.ErrIndexJobNotFound
	}

	return job, nil
}

func (r *Repository) Save(job *domain.IndexJob) error {

	r.jobs[job.PageID] = job

	return nil
}

func (r *Repository) ListByStatus(limit int, status domain.IndexStatus) ([]*domain.IndexJob, error) {
	var jobs []*domain.IndexJob

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

package mem

import "crawlquery/node/domain"

type Repository struct {
	jobs map[string]*domain.RepairJob
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make(map[string]*domain.RepairJob),
	}
}

func (r *Repository) Create(job *domain.RepairJob) error {
	r.jobs[job.PageID] = job
	return nil
}

func (r *Repository) Get(pageID string) (*domain.RepairJob, error) {
	job, ok := r.jobs[pageID]
	if !ok {
		return nil, domain.ErrRepairJobNotFound
	}
	return job, nil
}

func (r *Repository) Update(job *domain.RepairJob) error {
	r.jobs[job.PageID] = job
	return nil
}

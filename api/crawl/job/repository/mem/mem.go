package mem

import "crawlquery/api/domain"

type Repository struct {
	jobs map[string]*domain.CrawlJob
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make(map[string]*domain.CrawlJob),
	}
}

func (r *Repository) Create(j *domain.CrawlJob) error {
	r.jobs[j.ID] = j
	return nil
}

func (r *Repository) Get(id string) (*domain.CrawlJob, error) {
	job, ok := r.jobs[id]
	if !ok {
		return nil, nil
	}

	return job, nil
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

func (r *Repository) Delete(id string) error {
	delete(r.jobs, id)
	return nil
}

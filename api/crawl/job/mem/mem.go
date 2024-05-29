package mem

import "crawlquery/api/domain"

type Repository struct {
	jobs map[domain.PageID]*domain.CrawlJob
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make(map[domain.PageID]*domain.CrawlJob),
	}
}

func (r *Repository) Get(pageID domain.PageID) (*domain.CrawlJob, error) {
	job, ok := r.jobs[pageID]
	if !ok {
		return nil, domain.ErrCrawlJobNotFound
	}
	return job, nil
}

func (r *Repository) Save(cj *domain.CrawlJob) error {
	r.jobs[cj.PageID] = cj
	return nil
}

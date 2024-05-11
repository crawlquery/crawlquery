package mem

import "crawlquery/pkg/domain"

type MemoryRepository struct {
	jobs []*domain.CrawlJob
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (mr *MemoryRepository) Save() error {
	return nil
}

func (mr *MemoryRepository) Load() error {
	return nil
}

func (mr *MemoryRepository) Pop() (*domain.CrawlJob, error) {
	if len(mr.jobs) == 0 {
		return nil, nil
	}

	job := mr.jobs[0]
	mr.jobs = mr.jobs[1:]
	return job, nil
}

func (mr *MemoryRepository) Push(j *domain.CrawlJob) error {
	mr.jobs = append(mr.jobs, j)
	return nil
}

package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"time"
)

type Service struct {
	repo domain.CrawlJobRepository
}

func NewService(repo domain.CrawlJobRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (cs *Service) AddJob(url string) error {
	job := &domain.CrawlJob{
		ID:        util.UUID(),
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := job.Validate(); err != nil {
		return err
	}

	// Save the job in the repository
	if err := cs.repo.Create(job); err != nil {
		return err
	}
	return nil
}

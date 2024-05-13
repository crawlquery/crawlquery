package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo   domain.CrawlJobRepository
	logger *zap.SugaredLogger
}

func NewService(
	repo domain.CrawlJobRepository,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (cs *Service) Create(url string) (*domain.CrawlJob, error) {
	job := &domain.CrawlJob{
		ID:        util.UUID(),
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := job.Validate(); err != nil {
		return nil, err
	}

	// Save the job in the repository
	if err := cs.repo.Create(job); err != nil {
		cs.logger.Errorw("Crawl.Service.AddJob: error creating job", "error", err)
		return nil, domain.ErrInternalError
	}
	return job, nil
}
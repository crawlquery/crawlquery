package service

import (
	"crawlquery/api/domain"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	indexJobRepo domain.IndexJobRepository
	logger       *zap.SugaredLogger
}

func NewService(indexJobRepo domain.IndexJobRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		indexJobRepo: indexJobRepo,
		logger:       logger,
	}
}

func (s *Service) Create(pageID string) (*domain.IndexJob, error) {

	if _, err := s.indexJobRepo.GetByPageID(pageID); err == nil {
		return nil, domain.ErrIndexJobAlreadyExists
	}

	job := &domain.IndexJob{
		PageID:    pageID,
		CreatedAt: time.Now(),
	}

	return s.indexJobRepo.Create(job)
}

func (s *Service) Get(id string) (*domain.IndexJob, error) {
	return s.indexJobRepo.Get(id)
}

func (s *Service) Next() (*domain.IndexJob, error) {
	return s.indexJobRepo.Next()
}

func (s *Service) Update(job *domain.IndexJob) error {
	return s.indexJobRepo.Update(job)
}

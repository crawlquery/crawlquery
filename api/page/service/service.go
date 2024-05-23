package service

import (
	"crawlquery/api/domain"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	pageRepo domain.PageRepository
	logger   *zap.SugaredLogger
}

func NewService(pageRepo domain.PageRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		pageRepo: pageRepo,
		logger:   logger,
	}
}

func (s *Service) Get(pageID string) (*domain.Page, error) {
	page, err := s.pageRepo.Get(pageID)
	if err != nil {
		s.logger.Errorw("Error getting page", "error", err, "pageID", pageID)
		return nil, err
	}
	return page, nil
}

func (s *Service) Create(pageID string, shardID uint) (*domain.Page, error) {

	if _, err := s.pageRepo.Get(pageID); err == nil {
		s.logger.Errorw("Page already exists", "pageID", pageID)
		return nil, domain.ErrPageAlreadyExists
	}

	page := &domain.Page{
		ID:        pageID,
		ShardID:   shardID,
		CreatedAt: time.Now(),
	}

	if err := s.pageRepo.Create(page); err != nil {
		s.logger.Errorw("Error creating page", "error", err, "pageID", pageID)
		return nil, err
	}

	return page, nil
}

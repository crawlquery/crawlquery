package service

import (
	"crawlquery/api/domain"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	linkRepo domain.LinkRepository
	logger   *zap.SugaredLogger
}

func NewService(linkRepo domain.LinkRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		linkRepo: linkRepo,
		logger:   logger,
	}
}

func (s *Service) Create(srcID, dstID string) (*domain.Link, error) {
	link := &domain.Link{
		SrcID:     srcID,
		DstID:     dstID,
		CreatedAt: time.Now(),
	}

	err := s.linkRepo.Create(link)

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return link, nil
}

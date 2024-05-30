package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"time"
)

type Service struct {
	versionRepo domain.PageVersionRepository
}

type Option func(*Service)

func WithVersionRepo(versionRepo domain.PageVersionRepository) Option {
	return func(s *Service) {
		s.versionRepo = versionRepo
	}
}

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) Create(pageID domain.PageID, contentHash domain.ContentHash) (*domain.PageVersion, error) {
	pageVersion := &domain.PageVersion{
		ID:          domain.PageVersionID(util.UUIDString()),
		PageID:      pageID,
		ContentHash: contentHash,
		CreatedAt:   time.Now(),
	}

	err := s.versionRepo.Create(pageVersion)
	if err != nil {
		return nil, err
	}

	return pageVersion, nil
}

func (s *Service) Get(id domain.PageVersionID) (*domain.PageVersion, error) {
	pageVersion, err := s.versionRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return pageVersion, nil
}

func (s *Service) ListByPageID(pageID domain.PageID) ([]*domain.PageVersion, error) {
	versions, err := s.versionRepo.ListByPageID(pageID)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

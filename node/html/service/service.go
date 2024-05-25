package service

import "crawlquery/node/domain"

type Service struct {
	repo          domain.HTMLRepository
	backupService domain.HTMLBackupService
}

func NewService(repo domain.HTMLRepository, backupService domain.HTMLBackupService) *Service {
	return &Service{
		repo:          repo,
		backupService: backupService,
	}
}

func (s *Service) Save(pageID string, html []byte) error {
	err := s.repo.Save(pageID, html)

	if err != nil {
		return err
	}

	return s.backupService.Save(pageID, html)
}

func (s *Service) Get(pageID string) ([]byte, error) {
	data, err := s.repo.Get(pageID)

	if err == nil {
		return data, nil
	}

	return s.backupService.Get(pageID)
}

func (s *Service) Restore(pageID string) error {
	data, err := s.backupService.Get(pageID)

	if err != nil {
		return err
	}

	return s.repo.Save(pageID, data)
}

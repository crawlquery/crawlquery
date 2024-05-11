package service

import (
	"crawlquery/pkg/domain"
)

type IndexService struct {
	repo domain.IndexRepository
}

func NewIndexService(repo domain.IndexRepository) *IndexService {
	return &IndexService{
		repo: repo,
	}
}

func (service *IndexService) Search(query string) ([]domain.Result, error) {
	idx, err := service.repo.Load()
	if err != nil {
		return nil, err
	}

	return idx.Search(query)
}

func (service *IndexService) AddPage(p *domain.Page) error {
	idx, err := service.repo.Load()
	if err != nil {
		return err
	}

	idx.AddPage(p)

	return service.repo.Save(idx)
}

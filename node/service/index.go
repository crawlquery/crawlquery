package service

import (
	"crawlquery/pkg/domain"
)

type IndexService struct {
	repo  domain.IndexRepository
	index domain.Index
}

func NewIndexService(repo domain.IndexRepository) *IndexService {
	return &IndexService{
		repo: repo,
	}
}

func (service *IndexService) LoadIndex() error {
	idx, err := service.repo.Load()
	if err != nil {
		return err
	}
	service.index = idx
	return nil
}

func (service *IndexService) Search(query string) []domain.Result {
	return service.index.Search(query)
}

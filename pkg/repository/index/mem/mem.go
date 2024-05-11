package mem

import (
	"crawlquery/pkg/domain"
)

type MemoryRepository struct {
	index domain.Index
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (mr *MemoryRepository) Save(idx domain.Index) error {
	mr.index = idx
	return nil
}

func (mr *MemoryRepository) Load() (domain.Index, error) {
	return mr.index, nil
}

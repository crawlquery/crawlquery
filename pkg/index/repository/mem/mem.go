package mem

import "crawlquery/pkg/index"

type MemoryRepository struct {
	index *index.Index
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (mr *MemoryRepository) Save(idx *index.Index) error {
	mr.index = idx
	return nil
}

func (mr *MemoryRepository) Load() (*index.Index, error) {
	return mr.index, nil
}

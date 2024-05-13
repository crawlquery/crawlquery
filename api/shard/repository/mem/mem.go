package mem

import "crawlquery/api/domain"

type Repository struct {
	shards map[int]*domain.Shard
}

func NewRepository() *Repository {
	return &Repository{
		shards: make(map[int]*domain.Shard),
	}
}

func (r *Repository) Create(s *domain.Shard) error {
	r.shards[s.ID] = s
	return nil
}

func (r *Repository) List() ([]*domain.Shard, error) {
	shards := []*domain.Shard{}
	for _, s := range r.shards {
		shards = append(shards, s)
	}
	return shards, nil
}

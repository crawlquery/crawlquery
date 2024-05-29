package mem

import "crawlquery/api/domain"

type Repository struct {
	shards map[domain.ShardID]*domain.Shard
}

func NewRepository() *Repository {
	return &Repository{
		shards: make(map[domain.ShardID]*domain.Shard),
	}
}

func (r *Repository) Create(s *domain.Shard) error {
	r.shards[s.ID] = s
	return nil
}

func (r *Repository) Get(id domain.ShardID) (*domain.Shard, error) {
	s, ok := r.shards[id]
	if !ok {
		return nil, domain.ErrShardNotFound
	}
	return s, nil
}

func (r *Repository) List() ([]*domain.Shard, error) {
	shards := []*domain.Shard{}
	for _, s := range r.shards {
		shards = append(shards, s)
	}
	return shards, nil
}

func (r *Repository) Count() (int, error) {
	return len(r.shards), nil
}

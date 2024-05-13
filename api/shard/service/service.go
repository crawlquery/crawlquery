package service

import (
	"crawlquery/api/domain"
	"hash/fnv"

	"go.uber.org/zap"
)

type Service struct {
	repo   domain.ShardRepository
	logger *zap.SugaredLogger
}

func NewService(
	repo domain.ShardRepository,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{repo, logger}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (ss *Service) Create(s *domain.Shard) error {
	return ss.repo.Create(s)
}

func (ss *Service) First() (*domain.Shard, error) {
	shards, err := ss.repo.List()

	if err != nil {
		ss.logger.Errorf("Shard.Service.First: error listing shards: %v", err)
		return nil, err
	}

	if len(shards) == 0 {
		ss.logger.Errorf("Shard.Service.First: no shards")
		return nil, domain.ErrNoShards
	}

	return shards[0], nil
}

func (ss *Service) GetURLShardID(url string) (uint, error) {

	count, err := ss.repo.Count()

	if err != nil {
		ss.logger.Errorf("Shard.Service.GetURLShardID: error counting shards: %v", err)
		return 0, err
	}

	if count == 0 {
		ss.logger.Errorf("Shard.Service.GetURLShardID: no shards")
		return 0, domain.ErrNoShards
	}

	return uint(hash(url) % uint32(count)), nil
}

func (ss *Service) List() ([]*domain.Shard, error) {
	return ss.repo.List()
}

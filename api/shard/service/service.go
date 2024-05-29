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

type Option func(*Service)

func WithRepo(repo domain.ShardRepository) func(*Service) {
	return func(s *Service) {
		s.repo = repo
	}
}

func WithLogger(logger *zap.SugaredLogger) func(*Service) {
	return func(s *Service) {
		s.logger = logger
	}
}

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func hashURL(s domain.URL) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (ss *Service) GetURLShardID(url domain.URL) (domain.ShardID, error) {

	count, err := ss.repo.Count()

	if err != nil {
		ss.logger.Errorf("Shard.Service.GetURLShardID: error counting shards: %v", err)
		return 0, err
	}

	if count == 0 {
		ss.logger.Errorf("Shard.Service.GetURLShardID: no shards")
		return 0, domain.ErrNoShards
	}

	return domain.ShardID(uint(hashURL(url) % uint32(count))), nil
}

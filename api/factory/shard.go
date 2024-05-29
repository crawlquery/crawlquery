package factory

import (
	"crawlquery/api/domain"
	"crawlquery/api/shard/repository/mem"
	"crawlquery/api/shard/service"
	"crawlquery/pkg/testutil"
)

func ShardRepoWithShard(s *domain.Shard) *mem.Repository {
	repo := mem.NewRepository()

	if s != nil {
		repo.Create(s)
	}

	return repo
}

func ShardServiceWithShard(s *domain.Shard) (*service.Service, *mem.Repository) {
	repo := ShardRepoWithShard(s)

	return service.NewService(service.WithRepo(repo), service.WithLogger(
		testutil.NewTestLogger(),
	)), repo
}

package factory

import (
	"crawlquery/api/domain"
	"crawlquery/api/node/repository/mem"
	"crawlquery/api/node/service"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"

	"time"
)

func NodeRepoWithNode(n *domain.Node) *mem.Repository {
	repo := mem.NewRepository()

	if n != nil {
		if n.ID == "" {
			n.ID = util.UUIDString()
		}
		if n.CreatedAt.IsZero() {
			n.CreatedAt = time.Now()
		}

		if n.Hostname == "" {
			n.Hostname = "testnode"
		}

		repo.Create(n)
	}

	return repo
}

func NodeServiceWithNode(
	as domain.AccountService,
	n *domain.Node,
) (*service.Service, *mem.Repository) {
	repo := NodeRepoWithNode(n)
	return service.NewService(
		service.WithNodeRepo(repo),
		service.WithAccountService(as),
		service.WithLogger(testutil.NewTestLogger()),
		service.WithRandSeed(time.Now().Unix()),
	), repo
}

func NodeService(
	as domain.AccountService,
) (*service.Service, *mem.Repository) {
	return NodeServiceWithNode(as, nil)
}

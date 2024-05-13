package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo           domain.NodeRepository
	accountService domain.AccountService
	// shardService domain.ShardService
	logger *zap.SugaredLogger
}

func NewService(
	repo domain.NodeRepository,
	accountService domain.AccountService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		repo:           repo,
		accountService: accountService,
		logger:         logger,
	}
}

func (cs *Service) Create(accountID, hostname string, port uint) (*domain.Node, error) {

	// Check if the account exists
	_, err := cs.accountService.Get(accountID)

	if err != nil {
		return nil, domain.ErrAccountNotFound
	}

	node := &domain.Node{
		ID:        util.UUID(),
		AccountID: accountID,
		Hostname:  hostname,
		Port:      port,
		ShardID:   0,
		CreatedAt: time.Now(),
	}

	if err := node.Validate(); err != nil {
		return nil, err
	}

	// Save the node in the repository
	if err := cs.repo.Create(node); err != nil {
		cs.logger.Errorw("Node.Service.Create: error creating node", "error", err)
		return nil, domain.ErrInternalError
	}
	return node, nil
}

func (cs *Service) List() ([]*domain.Node, error) {
	return cs.repo.List()
}

func (cs *Service) RandomizedList() ([]*domain.Node, error) {
	nodes, err := cs.repo.List()
	if err != nil {
		return nil, err
	}

	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	return nodes, nil
}

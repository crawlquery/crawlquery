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

	// Check if the node already exists
	all, err := cs.repo.List()

	if err != nil {
		cs.logger.Errorf("Node.Service.Create: error listing nodes: %v", err)
		return nil, domain.ErrInternalError
	}

	for _, n := range all {
		if n.Hostname == hostname && n.Port == port {
			return nil, domain.ErrNodeAlreadyExists
		}
	}

	// Check if the account exists
	_, err = cs.accountService.Get(accountID)

	if err != nil {
		cs.logger.Errorf("Node.Service.Create: error getting account: %v", err)
		return nil, domain.ErrInvalidAccountID
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

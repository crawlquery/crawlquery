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
	shardService   domain.ShardService
	logger         *zap.SugaredLogger
}

func NewService(
	repo domain.NodeRepository,
	accountService domain.AccountService,
	shardService domain.ShardService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		repo:           repo,
		accountService: accountService,
		shardService:   shardService,
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
		CreatedAt: time.Now(),
	}

	if err := node.Validate(); err != nil {
		return nil, err
	}

	err = cs.AllocateNode(node)

	if err != nil {
		cs.logger.Errorw("Node.Service.Create: error allocating node", "error", err)
		return nil, domain.ErrInternalError
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

func (ss *Service) AllocateNode(node *domain.Node) error {
	// Get the shard with the least nodes
	shard, err := ss.GetShardWithLeastNodes()

	if err != nil {
		ss.logger.Errorf("Node.Allocation.Service.AllocateNode: error getting shard with least nodes: %v", err)

		shard, err = ss.shardService.First()

		if err != nil {
			ss.logger.Errorf("Node.Allocation.Service.AllocateNode: unable to create shard: %v", err)
			return domain.ErrInternalError
		}
	}

	node.ShardID = shard.ID

	return nil
}

func (ss *Service) GetShardWithLeastNodes() (*domain.Shard, error) {

	nodes, err := ss.List()

	if err != nil {
		ss.logger.Errorf("Shard.Service.GetShardWithLeastNodes: error listing nodes: %v", err)
		return nil, err
	}

	var shardDistribution = make(map[uint]int)

	for _, node := range nodes {
		shardDistribution[node.ShardID]++
	}

	var leastNodes *domain.Shard = nil
	var leastNodesCount int = int(^uint(0) >> 1) // Set to maximum int value

	for shardID, count := range shardDistribution {
		if leastNodes == nil || count < leastNodesCount {
			leastNodesCount = count
			// Assuming leastNodes needs to be newly initialized each time
			leastNodes = &domain.Shard{ID: shardID}
		}
	}

	if leastNodes == nil {
		ss.logger.Errorf("Shard.Service.GetShardWithLeastNodes: no shards")

		return nil, domain.ErrNoShards
	}

	return leastNodes, nil
}

func (s *Service) ListGroupByShard() (map[uint][]*domain.Node, error) {
	nodes, err := s.RandomizedList()
	if err != nil {
		return nil, err
	}

	grouped := make(map[uint][]*domain.Node)

	for _, node := range nodes {
		grouped[node.ShardID] = append(grouped[node.ShardID], node)
	}

	return grouped, nil
}

package service

import (
	"context"
	"crawlquery/api/domain"
	"crawlquery/node/dto"
	"crawlquery/pkg/client/node"
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

type Option func(*Service)

func WithNodeRepo(repo domain.NodeRepository) Option {
	return func(s *Service) {
		s.repo = repo
	}
}

func WithAccountService(accountService domain.AccountService) Option {
	return func(s *Service) {
		s.accountService = accountService
	}
}

func WithShardService(shardService domain.ShardService) Option {
	return func(s *Service) {
		s.shardService = shardService
	}
}

func WithLogger(logger *zap.SugaredLogger) Option {
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
		ID:        util.UUIDString(),
		AccountID: accountID,
		Key:       util.UUIDString(),
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
		return err
	}

	node.ShardID = shard.ID

	return nil
}

func (ss *Service) GetShardWithLeastNodes() (*domain.Shard, error) {

	shards, err := ss.shardService.List()

	if err != nil {
		ss.logger.Errorf("Shard.Service.GetShardWithLeastNodes: error listing shards: %v", err)
		return nil, err
	}

	nodes, err := ss.List()

	if err != nil {
		ss.logger.Errorf("Shard.Service.GetShardWithLeastNodes: error listing nodes: %v", err)
		return nil, err
	}

	var shardDistribution = make(map[*domain.Shard]int)

	for _, shard := range shards {
		shardDistribution[shard] = 0

		for _, node := range nodes {
			if node.ShardID == shard.ID {
				shardDistribution[shard]++
			}
		}
	}

	var minShard *domain.Shard

	for shard, count := range shardDistribution {
		if minShard == nil || count < shardDistribution[minShard] {
			minShard = shard
		}
	}

	if minShard == nil && len(shards) > 0 {
		return shards[0], nil
	}

	if minShard == nil {
		return nil, domain.ErrNoShards
	}

	return minShard, nil
}

func (s *Service) RandomizedListGroupByShard() (map[domain.ShardID][]*domain.Node, error) {
	nodes, err := s.RandomizedList()
	if err != nil {
		return nil, err
	}

	grouped := make(map[domain.ShardID][]*domain.Node)

	for _, node := range nodes {
		grouped[node.ShardID] = append(grouped[node.ShardID], node)
	}

	return grouped, nil
}

func (s *Service) ListByAccountID(accountID string) ([]*domain.Node, error) {
	return s.repo.ListByAccountID(accountID)
}

func (s *Service) Randomize(nodes []*domain.Node) []*domain.Node {
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	return nodes
}

func (s *Service) ListByShardID(shardID domain.ShardID) ([]*domain.Node, error) {
	nodes, err := s.List()
	if err != nil {
		return nil, err
	}

	var filtered []*domain.Node

	for _, node := range nodes {
		if node.ShardID == shardID {
			filtered = append(filtered, node)
		}
	}

	return filtered, nil
}

func (s *Service) SendCrawlJob(ctx context.Context, n *domain.Node, job *domain.CrawlJob) (*dto.CrawlResponse, error) {

	c := node.NewClient(
		node.WithHostname(n.Hostname),
		node.WithPort(n.Port),
		node.WithContext(ctx),
	)

	return c.Crawl(string(job.PageID), string(job.URL))
}

func (s *Service) SendIndexJob(n *domain.Node, job *domain.IndexJob) error {
	c := node.NewClient(
		node.WithHostname(n.Hostname),
		node.WithPort(n.Port),
	)

	return c.Index(job.PageID)
}

func (s *Service) Auth(key string) (*domain.Node, error) {
	node, err := s.repo.GetNodeByKey(key)

	if err != nil {
		return nil, domain.ErrNodeNotFound
	}

	return node, nil
}

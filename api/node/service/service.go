package service

import (
	"bytes"
	"crawlquery/api/domain"
	"crawlquery/pkg/dto"
	"crawlquery/pkg/util"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

func (s *Service) ListByAccountID(accountID string) ([]*domain.Node, error) {
	return s.repo.ListByAccountID(accountID)
}

func (s *Service) ListByShardID(shardID uint) ([]*domain.Node, error) {
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

func (s *Service) SendCrawlJob(node *domain.Node, job *domain.CrawlJob) error {

	req := dto.CrawlRequest{
		PageID: job.PageID,
		URL:    job.URL,
	}

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := http.Post(
		fmt.Sprintf("http://%s:%d/crawl", node.Hostname, node.Port),
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

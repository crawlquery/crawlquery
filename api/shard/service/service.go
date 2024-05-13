package service

import (
	"crawlquery/api/domain"
	"hash/fnv"

	"go.uber.org/zap"
)

type ShardService struct {
	repo        domain.ShardRepository
	nodeService domain.NodeService
	logger      *zap.SugaredLogger
}

func NewService(
	repo domain.ShardRepository,
	nodeService domain.NodeService,
	logger *zap.SugaredLogger,
) *ShardService {
	return &ShardService{repo, nodeService, logger}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (ss *ShardService) GetShardWithLeastNodes() (*domain.Shard, error) {

	nodes, err := ss.nodeService.List()

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

	return leastNodes, nil
}

func (ss *ShardService) GetURLShardID(url string) (int, error) {

	count, err := ss.repo.Count()

	if err != nil {
		ss.logger.Errorf("Shard.Service.GetURLShardID: error counting shards: %v", err)
		return 0, err
	}

	if count == 0 {
		ss.logger.Errorf("Shard.Service.GetURLShardID: no shards")
		return 0, domain.ErrNoShards
	}

	return int(hash(url) % uint32(count)), nil
}

package service

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/util"
	"errors"
	"math/rand"
	"net/url"
)

type NodeService struct {
	nr domain.NodeRepository
}

func NewNodeService(nr domain.NodeRepository) *NodeService {
	return &NodeService{
		nr: nr,
	}
}

func (service *NodeService) CreateOrUpdate(node *domain.Node) error {
	node.ID = util.UUID()
	return service.nr.CreateOrUpdate(node)
}

func (service *NodeService) Add(uri string) error {

	url, err := url.ParseRequestURI(uri)

	if err != nil {
		return err
	}

	hostname := url.Hostname()

	if hostname == "" {
		return errors.New("invalid url")
	}

	port := url.Port()

	if port == "" {
		port = "80"
	}

	return service.nr.CreateOrUpdate(&domain.Node{
		ID:       util.UUID(),
		Hostname: hostname,
		Port:     port,
		ShardID:  1,
	})
}

func (service *NodeService) Get(id string) (*domain.Node, error) {
	return service.nr.Get(id)
}

func (service *NodeService) GetRandom() (*domain.Node, error) {
	all, err := service.nr.GetAll()

	if err != nil {
		return nil, err
	}

	if len(all) == 0 {
		return nil, errors.New("no nodes found")
	}

	return all[rand.Intn(len(all))], nil
}

func (service *NodeService) RandomizeAll() ([]*domain.Node, error) {
	all, err := service.nr.GetAll()

	if err != nil {
		return nil, err
	}

	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})

	return all, nil
}

func (service *NodeService) AllByShard() (map[domain.ShardID][]*domain.Node, error) {
	all, err := service.RandomizeAll()

	if err != nil {
		return nil, err
	}

	shardNodes := map[domain.ShardID][]*domain.Node{}

	for _, node := range all {
		if shardNodes[node.ShardID] == nil {
			shardNodes[node.ShardID] = []*domain.Node{
				node,
			}
		} else {
			shardNodes[node.ShardID] = append(shardNodes[node.ShardID], node)
		}
	}

	return shardNodes, nil
}

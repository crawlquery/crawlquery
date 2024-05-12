package mem

import "crawlquery/pkg/domain"

type MemoryRepository struct {
	nodes map[string]*domain.Node
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nodes: make(map[string]*domain.Node),
	}
}

func (mr *MemoryRepository) CreateOrUpdate(n *domain.Node) error {
	mr.nodes[n.ID] = n
	return nil
}

func (mr *MemoryRepository) Create(n *domain.Node) error {
	mr.nodes[n.ID] = n
	return nil
}

func (mr *MemoryRepository) Get(id string) (*domain.Node, error) {
	return mr.nodes[id], nil
}

func (mr *MemoryRepository) GetAll() ([]*domain.Node, error) {
	nodes := make([]*domain.Node, 0, len(mr.nodes))
	for _, n := range mr.nodes {
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (mr *MemoryRepository) Delete(id string) error {
	delete(mr.nodes, id)
	return nil
}

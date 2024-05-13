package mem

import "crawlquery/api/domain"

type Repository struct {
	nodes            map[string]*domain.Node
	forceCreateError error
	forceListError   error
}

func NewRepository() *Repository {
	return &Repository{
		nodes: make(map[string]*domain.Node),
	}
}

func (r *Repository) ForceCreateError(e error) {
	r.forceCreateError = e
}

func (r *Repository) ForceListError(e error) {
	r.forceListError = e
}

func (r *Repository) Create(n *domain.Node) error {

	if r.forceCreateError != nil {
		return r.forceCreateError
	}
	r.nodes[n.ID] = n
	return nil
}

func (r *Repository) List() ([]*domain.Node, error) {
	if r.forceListError != nil {
		return nil, r.forceListError
	}
	nodes := make([]*domain.Node, 0)
	for _, n := range r.nodes {
		nodes = append(nodes, n)
	}

	return nodes, nil
}

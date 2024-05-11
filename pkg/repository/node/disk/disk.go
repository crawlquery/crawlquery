package disk

import (
	"crawlquery/pkg/domain"
	"encoding/gob"
	"errors"
	"os"
	"sync"
)

type DiskRepository struct {
	filepath string
	nodes    []*domain.Node
	lock     sync.Mutex
}

func NewDiskRepository(filepath string) *DiskRepository {
	return &DiskRepository{
		filepath: filepath,
	}
}

func (dr *DiskRepository) Load() error {
	file, err := os.Open(dr.filepath)

	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&dr.nodes); err != nil {
		return err
	}

	return nil
}

func (dr *DiskRepository) Save() error {
	dr.lock.Lock()
	defer dr.lock.Unlock()

	file, err := os.Create(dr.filepath)

	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(dr.nodes); err != nil {
		return err
	}

	return nil
}

func (dr *DiskRepository) CreateOrUpdate(n *domain.Node) error {
	for _, node := range dr.nodes {
		if node.ID == n.ID {
			// Update the node
			node = n
		}
	}

	dr.nodes = append(dr.nodes, n)
	dr.Save()

	return nil
}

func (dr *DiskRepository) Get(id string) (*domain.Node, error) {
	for _, node := range dr.nodes {
		if node.ID == id {
			return node, nil
		}
	}

	return nil, errors.New("node not found")
}

func (dr *DiskRepository) GetAll() ([]*domain.Node, error) {
	return dr.nodes, nil
}

func (dr *DiskRepository) Delete(id string) error {
	for i, node := range dr.nodes {
		if node.ID == id {
			dr.nodes = append(dr.nodes[:i], dr.nodes[i+1:]...)
			return nil
		}
	}
	dr.Save()

	return errors.New("node not found")
}

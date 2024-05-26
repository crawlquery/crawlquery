package mem

import (
	"errors"
	"time"
)

type Repository struct {
	ranks map[string]float64
}

func NewRepository() *Repository {
	return &Repository{
		ranks: make(map[string]float64),
	}
}

func (r *Repository) Update(keyword string, rank float64, createdAt time.Time) error {
	r.ranks[keyword] = rank

	return nil
}

func (r *Repository) Get(keyword string) (float64, error) {
	rank, ok := r.ranks[keyword]
	if !ok {
		return 0, errors.New("page rank not found")
	}
	return rank, nil
}

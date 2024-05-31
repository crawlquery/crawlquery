package mem

import (
	"crawlquery/api/domain"
	"errors"
	"time"
)

type Repository struct {
	ranks map[domain.PageID]float64
}

func NewRepository() *Repository {
	return &Repository{
		ranks: make(map[domain.PageID]float64),
	}
}

func (r *Repository) Update(pageID domain.PageID, rank float64, createdAt time.Time) error {
	r.ranks[pageID] = rank

	return nil
}

func (r *Repository) Get(pageID domain.PageID) (float64, error) {
	rank, ok := r.ranks[pageID]
	if !ok {
		return 0, errors.New("page rank not found")
	}
	return rank, nil
}

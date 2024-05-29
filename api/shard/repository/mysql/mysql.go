package mysql

import (
	"crawlquery/api/domain"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(shard *domain.Shard) error {
	_, err := r.db.Exec("INSERT INTO shards (id, created_at) VALUES (?, ?)", shard.ID, shard.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating shard: %w", err)
	}

	return nil
}

func (r *Repository) Get(id domain.ShardID) (*domain.Shard, error) {
	row := r.db.QueryRow("SELECT id, created_at FROM shards WHERE id = ?", id)

	shard := &domain.Shard{}

	err := row.Scan(&shard.ID, &shard.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("error getting shard: %w", err)
	}

	return shard, nil
}

func (r *Repository) List() ([]*domain.Shard, error) {
	rows, err := r.db.Query("SELECT id, created_at FROM shards")
	if err != nil {
		return nil, fmt.Errorf("error listing shards: %w", err)
	}
	defer rows.Close()

	shards := []*domain.Shard{}
	for rows.Next() {
		shard := &domain.Shard{}
		err := rows.Scan(&shard.ID, &shard.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning shard: %w", err)
		}
		shards = append(shards, shard)
	}

	return shards, nil
}

func (r *Repository) Count() (int, error) {
	row := r.db.QueryRow("SELECT COUNT(id) FROM shards")

	var count int
	err := row.Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("error counting shards: %w", err)
	}

	return count, nil
}

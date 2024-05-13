package mysql

import (
	"crawlquery/api/domain"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(n *domain.Node) error {
	_, err := r.db.Exec("INSERT INTO nodes (id, account_id, hostname, port, shard_id, created_at) VALUES (?, ?, ?, ?, ?, ?)", n.ID, n.AccountID, n.Hostname, n.Port, n.ShardID, time.Now())
	return err
}

func (r *Repository) List() ([]*domain.Node, error) {
	rows, err := r.db.Query("SELECT id, account_id, hostname, port, shard_id, created_at FROM nodes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make([]*domain.Node, 0)
	for rows.Next() {
		var n domain.Node
		err := rows.Scan(&n.ID, &n.AccountID, &n.Hostname, &n.Port, &n.ShardID, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}

	return nodes, nil
}

func (r *Repository) ListByAccountID(accountID string) ([]*domain.Node, error) {
	rows, err := r.db.Query("SELECT id, account_id, hostname, port, shard_id, created_at FROM nodes WHERE account_id = ?", accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make([]*domain.Node, 0)
	for rows.Next() {
		var n domain.Node
		err := rows.Scan(&n.ID, &n.AccountID, &n.Hostname, &n.Port, &n.ShardID, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}

	return nodes, nil
}

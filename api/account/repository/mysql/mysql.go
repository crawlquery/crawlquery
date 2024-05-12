package mysql

import (
	"crawlquery/api/domain"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(a *domain.Account) error {
	_, err := r.db.Exec("INSERT INTO accounts (id, email, password, created_at) VALUES (?, ?, ?, ?)", a.ID, a.Email, a.Password, a.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(id string) (*domain.Account, error) {
	row := r.db.QueryRow("SELECT id, email, password, created_at FROM accounts WHERE id = ?", id)

	var account domain.Account
	err := row.Scan(&account.ID, &account.Email, &account.Password, &account.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNoAccountFound
	}

	return &account, err
}

func (r *Repository) GetByEmail(email string) (*domain.Account, error) {
	row := r.db.QueryRow("SELECT id, email, password, created_at FROM accounts WHERE email = ?", email)

	var account domain.Account
	err := row.Scan(&account.ID, &account.Email, &account.Password, &account.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNoAccountFound
	}

	return &account, err
}

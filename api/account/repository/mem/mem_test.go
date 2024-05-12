package mem

import (
	"crawlquery/api/domain"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	repo := NewRepository()

	account := &domain.Account{
		ID:        "account1",
		Email:     "test@example.com",
		Password:  "password",
		CreatedAt: time.Now().UTC(),
	}

	err := repo.Create(account)

	if err != nil {
		t.Fatalf("Error creating account: %v", err)
	}

	if repo.accounts[account.ID].ID != account.ID {
		t.Errorf("Expected ID to be %s, got %s", account.ID, repo.accounts[account.ID].ID)
	}

	if repo.accounts[account.ID].Email != account.Email {
		t.Errorf("Expected Email to be %s, got %s", account.Email, repo.accounts[account.ID].Email)
	}

	if repo.accounts[account.ID].Password != account.Password {
		t.Errorf("Expected Password to be %s, got %s", account.Password, repo.accounts[account.ID].Password)
	}

	if repo.accounts[account.ID].CreatedAt.Sub(account.CreatedAt) > time.Second || account.CreatedAt.Sub(repo.accounts[account.ID].CreatedAt) > time.Second {
		t.Errorf("Expected CreatedAt to be within one second of %v, got %v", account.CreatedAt, repo.accounts[account.ID].CreatedAt)
	}
}

func TestGet(t *testing.T) {
	repo := NewRepository()

	account := &domain.Account{
		ID:        "account1",
		Email:     "test@example.com",
		Password:  "password",
		CreatedAt: time.Now().UTC(),
	}

	repo.accounts[account.ID] = account

	got, err := repo.Get(account.ID)

	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}

	if got.ID != account.ID {
		t.Errorf("Expected ID to be %s, got %s", account.ID, got.ID)
	}

	if got.Email != account.Email {
		t.Errorf("Expected Email to be %s, got %s", account.Email, got.Email)
	}

	if got.Password != account.Password {
		t.Errorf("Expected Password to be %s, got %s", account.Password, got.Password)
	}

	if got.CreatedAt.Sub(account.CreatedAt) > time.Second || account.CreatedAt.Sub(got.CreatedAt) > time.Second {
		t.Errorf("Expected CreatedAt to be within one second of %v, got %v", account.CreatedAt, got.CreatedAt)
	}
}

func TestGetByEmail(t *testing.T) {
	repo := NewRepository()

	account := &domain.Account{
		ID:        "account1",
		Email:     "test@example.com",
		Password:  "password",
		CreatedAt: time.Now().UTC(),
	}

	repo.accounts[account.ID] = account

	got, err := repo.GetByEmail(account.Email)

	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}

	if got.ID != account.ID {
		t.Errorf("Expected ID to be %s, got %s", account.ID, got.ID)
	}

	if got.Email != account.Email {
		t.Errorf("Expected Email to be %s, got %s", account.Email, got.Email)
	}

	if got.Password != account.Password {
		t.Errorf("Expected Password to be %s, got %s", account.Password, got.Password)
	}

	if got.CreatedAt.Sub(account.CreatedAt) > time.Second || account.CreatedAt.Sub(got.CreatedAt) > time.Second {
		t.Errorf("Expected CreatedAt to be within one second of %v, got %v", account.CreatedAt, got.CreatedAt)
	}
}

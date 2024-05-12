package mem

import (
	"crawlquery/api/domain"
	"errors"
)

type Repository struct {
	accounts             map[string]*domain.Account
	forceCreateError     error
	forceGetByEmailError error
}

func NewRepository() *Repository {
	return &Repository{
		accounts: make(map[string]*domain.Account),
	}
}

func (r *Repository) ForceCreateError(err error) {
	r.forceCreateError = err
}

func (r *Repository) ForceGetByEmailError(err error) {
	r.forceGetByEmailError = err
}

func (r *Repository) Create(a *domain.Account) error {

	if r.forceCreateError != nil {
		return r.forceCreateError
	}

	if _, ok := r.accounts[a.ID]; ok {
		return errors.New("account already exists")
	}

	r.accounts[a.ID] = a
	return nil
}

func (r *Repository) Get(id string) (*domain.Account, error) {
	account, ok := r.accounts[id]
	if !ok {
		return nil, nil
	}

	return account, nil
}

func (r *Repository) GetByEmail(email string) (*domain.Account, error) {
	if r.forceGetByEmailError != nil {
		return nil, r.forceGetByEmailError
	}
	for _, account := range r.accounts {
		if account.Email == email {
			return account, nil
		}
	}

	return nil, nil
}

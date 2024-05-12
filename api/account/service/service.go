package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"time"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	repo      domain.AccountRepository
	validator *validator.Validate
}

func NewService(repo domain.AccountRepository) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *Service) Create(email, password string) (*domain.Account, error) {

	a := &domain.Account{
		ID:        util.UUID(),
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	if err := a.Validate(); err != nil {
		return nil, err
	}

	check, err := s.repo.GetByEmail(email)

	if err != nil {
		return nil, domain.InternalError
	}

	if check != nil {
		return nil, domain.ErrAccountExists
	}

	err = s.repo.Create(a)

	if err != nil {
		return nil, domain.InternalError
	}

	return a, nil
}

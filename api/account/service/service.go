package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Service struct {
	repo      domain.AccountRepository
	logger    *zap.SugaredLogger
	validator *validator.Validate
}

func NewService(
	repo domain.AccountRepository,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
		logger:    logger,
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
	if err == nil || check != nil {
		return nil, domain.ErrAccountExists
	}

	err = s.repo.Create(a)

	if err != nil {
		s.logger.Errorw("Account.Service.Create: error creating account", "error", err)
		return nil, domain.ErrInternalError
	}

	return a, nil
}

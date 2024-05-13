package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/authutil"
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

	if len(password) < 6 {
		return nil, domain.ErrPasswordTooShort
	}

	hashedPassword, err := authutil.HashPassword(password)
	if err != nil {
		s.logger.Errorw("Account.Service.Create: error hashing password", "error", err)
		return nil, domain.ErrInternalError
	}

	a := &domain.Account{
		ID:        util.UUID(),
		Email:     email,
		Password:  hashedPassword,
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

func (s *Service) Get(id string) (*domain.Account, error) {
	return s.repo.Get(id)
}

func (s *Service) GetByEmail(email string) (*domain.Account, error) {
	return s.repo.GetByEmail(email)
}

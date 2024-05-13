package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/authutil"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	accountService domain.AccountService
	logger         *zap.SugaredLogger
}

func NewService(
	accountService domain.AccountService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		accountService: accountService,
		logger:         logger,
	}
}

func (s *Service) Login(email, password string) (string, error) {
	acc, err := s.accountService.GetByEmail(email)

	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	token, err := authutil.GenerateToken(acc.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

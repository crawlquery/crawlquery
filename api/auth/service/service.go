package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/authutil"

	"go.uber.org/zap"
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
		if err != domain.ErrAccountNotFound {
			s.logger.Errorw("Auth.Service.Login: error getting account", "error", err)
		}
		return "", domain.ErrInvalidCredentials
	}

	err = authutil.CompareHashAndPassword(acc.Password, password)

	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := authutil.GenerateToken(acc.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

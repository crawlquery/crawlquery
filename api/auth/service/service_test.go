package service_test

import (
	"crawlquery/api/auth/service"
	"crawlquery/api/domain"
	"crawlquery/api/factory"
	"crawlquery/pkg/authutil"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("can log in", func(t *testing.T) {
		hashedPassword, err := authutil.HashPassword("password")
		if err != nil {
			t.Fatalf("Error hashing password: %v", err)
		}

		account := &domain.Account{
			ID:       util.UUID(),
			Email:    "test@example.com",
			Password: hashedPassword,
		}
		accSvc, _ := factory.AccountServiceWithAccount(
			account,
		)

		authSvc := service.NewService(accSvc, testutil.NewTestLogger())

		token, err := authSvc.Login(account.Email, "password")

		if err != nil {
			t.Fatalf("Error logging in: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be set, got empty string")
		}

		claims, err := authutil.ParseClaims(token)

		if err != nil {
			t.Fatalf("Error parsing token: %v", err)
		}

		id, ok := claims["id"].(string)

		if !ok {
			t.Errorf("Expected id to be set in claims")
		}

		if id != account.ID {
			t.Errorf("Expected id to be %s, got %s", account.ID, id)
		}

	})

	t.Run("can't log in with wrong password", func(t *testing.T) {
		hashedPassword, err := authutil.HashPassword("password")
		if err != nil {
			t.Fatalf("Error hashing password: %v", err)
		}

		account := &domain.Account{
			ID:       util.UUID(),
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		accSvc, _ := factory.AccountServiceWithAccount(
			account,
		)

		authSvc := service.NewService(accSvc, testutil.NewTestLogger())

		token, err := authSvc.Login(account.Email, "wrongpassword")

		if err != domain.ErrInvalidCredentials {
			t.Errorf("Expected error logging in with wrong password")
		}

		if token != "" {
			t.Errorf("Expected token to be empty, got %s", token)
		}
	})
}

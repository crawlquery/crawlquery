package service_test

import (
	"crawlquery/api/account/repository/mem"
	"crawlquery/api/account/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/testutil"
	"errors"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create an account", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		account := &domain.Account{
			Email:    "test@example.com",
			Password: "password",
		}

		account, err := svc.Create(account.Email, account.Password)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		check, err := repo.Get(account.ID)

		if err != nil {
			t.Fatalf("Error getting account: %v", err)
		}

		if check.Email != account.Email {
			t.Errorf("Expected Email to be %s, got %s", account.Email, check.Email)
		}

		if check.Password != account.Password {
			t.Errorf("Expected Password to be %s, got %s", account.Password, check.Password)
		}

		if check.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set, got zero value")
		}
	})

	t.Run("can't create an account with the same email", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		email := "test@example.com"
		password := "password"
		account, err := svc.Create(email, password)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		_, err = svc.Create(email, password+"1")

		if err == nil {
			t.Errorf("Expected error creating account with the same email")
		}

		_, err = repo.Get(account.ID)

		if err != nil {
			t.Fatalf("Error getting account: %v", err)
		}

		if account.Password != password {
			t.Errorf("Expected Password to be %s, got %s", password, account.Password)
		}
	})

	t.Run("validates account", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		t.Run("invalid email", func(t *testing.T) {
			_, err := svc.Create("invalid", "password")

			if err == nil {
				t.Errorf("Expected account to be invalid, got nil")
			}
		})

		t.Run("invalid password", func(t *testing.T) {
			_, err := svc.Create("test@example.com", "pass")

			if err == nil {
				t.Errorf("Expected account to be invalid, got nil")
			}
		})
	})

	t.Run("handles create repository error", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		expectErr := errors.New("db locked")
		repo.ForceCreateError(expectErr)

		_, err := svc.Create("test@example.com", "password")

		if err != domain.InternalError {
			t.Errorf("Expected error creating account, got %v", err)
		}
	})

	t.Run("handles email unique check error", func(t *testing.T) {
		email := "test@example.com"
		repo := mem.NewRepository()
		repo.Create(&domain.Account{
			Email: email,
		})
		svc := service.NewService(repo, testutil.NewTestLogger())

		_, err := svc.Create(email, "password")

		if err != domain.ErrAccountExists {
			t.Errorf("Expected error creating account, got %v", err)
		}
	})

}

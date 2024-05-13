package service_test

import (
	"crawlquery/api/account/repository/mem"
	"crawlquery/api/account/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/authutil"
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

		if err := authutil.CompareHashAndPassword(check.Password, "password"); err != nil {
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
		repo.Create(&domain.Account{
			Email:    email,
			Password: password,
		})

		_, err := svc.Create(email, password)

		if err != domain.ErrAccountExists {
			t.Errorf("Expected error creating account with the same email")
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

		if err != domain.ErrInternalError {
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

func TestGet(t *testing.T) {
	t.Run("can get an account", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		account := &domain.Account{
			Email:    "test@example.com",
			Password: "password",
		}

		repo.Create(account)

		check, err := svc.Get(account.ID)

		if err != nil {
			t.Fatalf("Error getting account: %v", err)
		}

		if check.Email != account.Email {
			t.Errorf("Expected Email to be %s, got %s", account.Email, check.Email)
		}

		if check.Password != account.Password {
			t.Errorf("Expected Password to be %s, got %s", account.Password, check.Password)
		}
	})

	t.Run("returns error if account not found", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		_, err := svc.Get("1")

		if err != domain.ErrAccountNotFound {
			t.Errorf("Expected error getting account, got %v", err)
		}
	})
}

func TestGetByEmail(t *testing.T) {
	t.Run("can get an account by email", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		account := &domain.Account{
			Email:    "test@example.com",
			Password: "password",
		}

		repo.Create(account)

		check, err := svc.GetByEmail(account.Email)

		if err != nil {
			t.Fatalf("Error getting account: %v", err)
		}

		if check.Email != account.Email {
			t.Errorf("Expected Email to be %s, got %s", account.Email, check.Email)
		}

		if check.Password != account.Password {
			t.Errorf("Expected Password to be %s, got %s", account.Password, check.Password)
		}
	})

	t.Run("returns error if account not found", func(t *testing.T) {
		repo := mem.NewRepository()
		svc := service.NewService(repo, testutil.NewTestLogger())

		_, err := svc.GetByEmail("aa")

		if err != domain.ErrAccountNotFound {
			t.Errorf("Expected error getting account, got %v", err)
		}

	})

}

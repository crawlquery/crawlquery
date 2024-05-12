package domain_test

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"strings"
	"testing"
	"time"
)

func TestAccountValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now(),
		}

		err := a.Validate()

		if err != nil {
			t.Errorf("Expected account to be valid, got error: %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		a := &domain.Account{
			ID:        "abs",
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now(),
		}

		err := a.Validate()

		if err == nil {
			t.Errorf("Expected account to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Account.ID") {
			t.Errorf("Expected error to contain 'Account.ID', got %v", err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		a := &domain.Account{
			ID:        util.UUID(),
			Email:     "test",
			Password:  "password",
			CreatedAt: time.Now(),
		}

		err := a.Validate()

		if err == nil {
			t.Errorf("Expected account to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Account.Email") {
			t.Errorf("Expected error to contain 'Account.Email', got %v", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		a := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "pass",
			CreatedAt: time.Now(),
		}

		err := a.Validate()

		if err == nil {
			t.Errorf("Expected account to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Account.Password") {
			t.Errorf("Expected error to contain 'Account.Password', got %v", err)
		}
	})

	t.Run("invalid created at", func(t *testing.T) {
		a := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Time{},
		}

		err := a.Validate()

		if err == nil {
			t.Errorf("Expected account to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Account.CreatedAt") {
			t.Errorf("Expected error to contain 'Account.CreatedAt', got %v", err)
		}
	})
}

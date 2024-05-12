package dto_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"testing"
	"time"
)

func TestNewCreateAccountResponse(t *testing.T) {
	t.Run("should return correct CreateAccountResponse from Account", func(t *testing.T) {
		// given
		a := &domain.Account{
			ID:        "1",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		}

		r := dto.NewCreateAccountResponse(a)

		// then
		if r.Account.ID != a.ID {
			t.Errorf("Expected ID to be %s, got %s", a.ID, r.Account.ID)
		}

		if r.Account.Email != a.Email {
			t.Errorf("Expected Email to be %s, got %s", a.Email, r.Account.Email)
		}

		if r.Account.CreatedAt != a.CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, got %v", a.CreatedAt, r.Account.CreatedAt)
		}

	})
}

func TestCreateAccountRequestToJson(t *testing.T) {
	t.Run("should return correct JSON from CreateAccountRequest", func(t *testing.T) {
		// given
		r := &dto.CreateAccountRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		// when
		b, err := r.ToJson()

		// then
		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		expected := `{"email":"test@example.com","password":"password"}`

		if string(b) != expected {
			t.Errorf("Expected JSON to be %s, got %s", expected, string(b))
		}
	})
}

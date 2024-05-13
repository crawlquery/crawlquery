package dto_test

import (
	"crawlquery/api/dto"
	"testing"
)

func TestNewLoginResponse(t *testing.T) {
	t.Run("can create a new login response", func(t *testing.T) {
		token := "token"

		resp := dto.NewLoginResponse(token)

		if resp.Token != token {
			t.Errorf("Expected token to be %s, got %s", token, resp.Token)
		}

	})
}

func TestLoginRequestToJSON(t *testing.T) {
	t.Run("can marshal login request to JSON", func(t *testing.T) {
		req := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		b, err := req.ToJSON()

		if err != nil {
			t.Fatalf("Error marshalling to JSON: %v", err)
		}

		if string(b) != `{"email":"test@example.com","password":"password"}` {
			t.Errorf("Expected JSON to be %s, got %s", `{"email":"test@example.com","password":"password"}`, string(b))
		}
	})
}

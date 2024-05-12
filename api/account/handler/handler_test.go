package handler_test

import (
	"bytes"
	"crawlquery/api/account/handler"
	"crawlquery/api/account/repository/mem"
	"crawlquery/api/account/service"
	"crawlquery/api/dto"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreate(t *testing.T) {
	t.Run("should create an account", func(t *testing.T) {

		repo := mem.NewRepository()
		svc := service.NewService(repo)
		handler := handler.NewAccountHandler(svc)

		// given
		a := &dto.CreateAccountRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		req, err := a.ToJson()

		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer(req))

		// when
		handler.Create(ctx)

		// then
		if ctx.Writer.Status() != http.StatusCreated {
			t.Errorf("Expected status to be 201, got %d", ctx.Writer.Status())
		}

		var res dto.CreateAccountResponse
		err = json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if res.Account.Email != a.Email {
			t.Errorf("Expected email to be %s, got %s", a.Email, res.Account.Email)
		}

		if res.Account.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be non-zero")
		}
	})

	t.Run("should return 400 if malformed JSON", func(t *testing.T) {

		repo := mem.NewRepository()
		svc := service.NewService(repo)
		handler := handler.NewAccountHandler(svc)

		// given
		req := []byte(`{"email":`)
		responseWriter := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer(req))

		// when
		handler.Create(ctx)

		// then
		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}

		var res dto.ErrorResponse
		err := json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if res.Error == "" {
			t.Errorf("Expected error message, got empty string")
		}

	})

	t.Run("should return 400 if request data is invalid", func(t *testing.T) {

		repo := mem.NewRepository()
		svc := service.NewService(repo)
		handler := handler.NewAccountHandler(svc)

		// given
		a := &dto.CreateAccountRequest{
			Email:    "invalidemail",
			Password: "password",
		}

		req, err := a.ToJson()

		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		responseWriter := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer(req))

		// when
		handler.Create(ctx)

		// then
		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}

		var res dto.ErrorResponse
		err = json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if res.Error == "" {
			t.Errorf("Expected error message, got empty string")
		}

	})
}

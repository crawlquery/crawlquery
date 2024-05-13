package handler_test

import (
	"bytes"
	"crawlquery/api/auth/handler"
	"crawlquery/api/auth/service"
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/factory"
	"crawlquery/pkg/authutil"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	t.Run("should login", func(t *testing.T) {

		hashedPassword, err := authutil.HashPassword("password")

		if err != nil {
			t.Fatalf("Error hashing password: %v", err)
		}

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID:       accountID,
			Email:    "test@example.com",
			Password: hashedPassword,
		})

		svc := service.NewService(accSvc, testutil.NewTestLogger())
		handler := handler.NewHandler(svc)

		// given
		a := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		req, err := a.ToJSON()

		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(req))

		// when
		handler.Login(ctx)

		// then
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status to be 200, got %d", ctx.Writer.Status())
		}

		var res dto.LoginResponse
		err = json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if res.Token == "" {
			t.Error("Expected token to be set, got empty string")
		}

		claims, err := authutil.ParseClaims(res.Token)

		if err != nil {
			t.Fatalf("Error parsing token: %v", err)
		}

		id, ok := claims["id"].(string)

		if !ok {
			t.Errorf("Expected id to be set in claims")
		}

		if id != accountID {
			t.Errorf("Expected id to be %s, got %s", accountID, id)
		}
	})

	t.Run("should return 400 if malformed JSON", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc := service.NewService(accSvc, testutil.NewTestLogger())
		handler := handler.NewHandler(svc)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("{")))

		handler.Login(ctx)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}
	})

	t.Run("should return 401 if email is empty", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc := service.NewService(accSvc, testutil.NewTestLogger())
		handler := handler.NewHandler(svc)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"password":"password"}`)))

		handler.Login(ctx)

		if ctx.Writer.Status() != http.StatusUnauthorized {
			t.Errorf("Expected status to be 401, got %d", ctx.Writer.Status())
		}
	})
}

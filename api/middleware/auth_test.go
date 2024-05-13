package middleware_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/factory"
	"crawlquery/api/middleware"
	"crawlquery/pkg/authutil"
	"crawlquery/pkg/util"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "secret")

		account := &domain.Account{
			ID:    util.UUID(),
			Email: "test@example.com",
		}

		token, err := authutil.GenerateToken(account.ID)

		if err != nil {
			t.Fatalf("Error generating token: %v", err)
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		next := func(c *gin.Context) {
			acc, ok := c.Get("account")
			if !ok {
				t.Errorf("Expected account to be set in context")
			}

			if acc.(*domain.Account).ID != account.ID {
				t.Errorf("Expected account ID to be %s, got %s", account.ID, acc.(*domain.Account).ID)
			}

			c.JSON(200, nil)
		}

		svc, _ := factory.AccountServiceWithAccount(account)

		mw := middleware.AuthMiddleware(svc, next)

		c, _ := gin.CreateTestContext(rec)

		c.Request = req
		mw(c)

		if rec.Code != 200 {
			t.Errorf("Expected status code to be 200, got %d", rec.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "secret")

		req := httptest.NewRequest("GET", "/", nil)

		rec := httptest.NewRecorder()

		next := func(c *gin.Context) {
			t.Errorf("Expected next to not be called")
		}

		svc, _ := factory.AccountServiceWithAccount(nil)

		mw := middleware.AuthMiddleware(svc, next)

		c, _ := gin.CreateTestContext(rec)

		c.Request = req
		mw(c)

		if rec.Code != 401 {
			t.Errorf("Expected status code to be 401, got %d", rec.Code)
		}
	})

}

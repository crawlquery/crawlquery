package errorutil_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/errorutil"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleGinError(t *testing.T) {
	t.Run("send internal error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

		errorutil.HandleGinError(ctx, domain.ErrInternalError, http.StatusBadRequest)

		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, ctx.Writer.Status())
		}
	})

	t.Run("send default error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

		errorutil.HandleGinError(ctx, errors.New("some error"), http.StatusBadRequest)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, ctx.Writer.Status())
		}
	})

	t.Run("do nothing if no error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

		errorutil.HandleGinError(ctx, nil, http.StatusBadRequest)

		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, ctx.Writer.Status())
		}
	})
}

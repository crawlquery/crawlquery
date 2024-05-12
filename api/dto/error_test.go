package dto_test

import (
	"crawlquery/api/dto"
	"fmt"
	"testing"
)

func TestNewErrorResponse(t *testing.T) {
	t.Run("should return correct ErrorResponse from error", func(t *testing.T) {
		// given
		err := fmt.Errorf("test error")

		r := dto.NewErrorResponse(err)

		// then
		if r.Error != err.Error() {
			t.Errorf("Expected Error to be %s, got %s", err.Error(), r.Error)
		}
	})
}

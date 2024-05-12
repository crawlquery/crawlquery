package errorutil

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGinError(c *gin.Context, err error, defaultCode int) {
	if err == nil {
		return
	}

	switch err {
	case domain.ErrInternalError:
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(err))
	default:
		c.JSON(defaultCode, dto.NewErrorResponse(err))
	}
}

package middleware

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/pkg/authutil"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(as domain.AccountService, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.GetHeader("Authorization")

		if jwtToken == "" {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(errors.New("unauthorized")))
			c.Abort()
			return
		}

		claims, err := authutil.ParseClaims(jwtToken[7:])

		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(errors.New("unauthorized")))
			c.Abort()
			return
		}

		id, ok := claims["id"].(string)

		if !ok {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(errors.New("unauthorized")))
			c.Abort()
			return
		}

		account, err := as.Get(id)

		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(errors.New("unauthorized")))
			c.Abort()
			return
		}
		c.Set("account", account)
		next(c)
	}
}

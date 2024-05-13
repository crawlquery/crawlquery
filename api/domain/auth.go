package domain

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService interface {
	Login(email, password string) (string, error)
}

type AuthHandler interface {
	Login(c *gin.Context)
}

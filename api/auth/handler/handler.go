package handler

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/errorutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorResponse(err))
		return
	}

	token, err := ah.authService.Login(req.Email, req.Password)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, dto.NewLoginResponse(token))
}

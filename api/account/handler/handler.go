package handler

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/errorutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService domain.AccountService
}

func NewHandler(accountService domain.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (ah *AccountHandler) Create(c *gin.Context) {
	var req dto.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.NewErrorResponse(err))
		return
	}

	account, err := ah.accountService.Create(req.Email, req.Password)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated,
		dto.NewCreateAccountResponse(account),
	)
}

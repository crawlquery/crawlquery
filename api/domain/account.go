package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrAccountExists = errors.New("cannot create account")
var ErrAccountNotFound = errors.New("account not found")

type Account struct {
	ID        string    `json:"id" validate:"required,uuid"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-" validate:"required,min=6,max=100"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

func (a *Account) Validate() error {
	return validate.Struct(a)
}

type AccountRepository interface {
	Create(*Account) error
	Get(string) (*Account, error)
	GetByEmail(string) (*Account, error)
}

type AccountService interface {
	Create(email, password string) (*Account, error)
	Get(string) (*Account, error)
	GetByEmail(string) (*Account, error)
}

type AccountHandler interface {
	Create(c *gin.Context)
}

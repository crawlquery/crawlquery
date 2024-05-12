package domain

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrAccountExists = errors.New("account already exists")

type Account struct {
	ID        string    `validate:"required,uuid"`
	Email     string    `validate:"required,email"`
	Password  string    `validate:"required,min=6,max=100"`
	CreatedAt time.Time `validate:"required"`
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
}

type AccountHandler interface {
	Create(c *gin.Context)
}

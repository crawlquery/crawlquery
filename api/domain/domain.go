package domain

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(
	validator.WithRequiredStructEnabled(),
)

var InternalError = errors.New("internal error")

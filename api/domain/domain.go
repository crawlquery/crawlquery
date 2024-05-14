package domain

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(
	validator.WithRequiredStructEnabled(),
)

var ErrInternalError = errors.New("internal error")
var ErrForbidden = errors.New("forbidden")
var ErrNoRowsUpdated = errors.New("no rows updated")

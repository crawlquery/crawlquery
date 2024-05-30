package domain

import "errors"

var ErrHTMLNotFound = errors.New("html not found")

type HTMLService interface {
	Get(hash string) ([]byte, error)
	Save(hash string, html []byte) error
}

type HTMLRepository interface {
	Get(hash string) ([]byte, error)
	Save(hash string, html []byte) error
}

type HTMLBackupService interface {
	Get(hash string) ([]byte, error)
	Save(hash string, html []byte) error
}

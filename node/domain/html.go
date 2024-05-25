package domain

import "errors"

var ErrHTMLNotFound = errors.New("html not found")

type HTMLService interface {
	Get(pageID string) ([]byte, error)
	Save(pageID string, html []byte) error
}

type HTMLRepository interface {
	Get(pageID string) ([]byte, error)
	Save(pageID string, html []byte) error
}

type HTMLBackupService interface {
	Get(pageID string) ([]byte, error)
	Save(pageID string, html []byte) error
}

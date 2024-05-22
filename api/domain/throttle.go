package domain

import (
	"database/sql"
	"errors"
)

var ErrInvalidLockKey = errors.New("invalid lock key")
var ErrDomainLocked = errors.New("domain is locked")
var ErrDomainNotLocked = errors.New("domain is not locked")

type Lock struct {
	Domain   string
	Key      string
	LockedAt sql.NullTime
}

type DomainLockRepository interface {
	IsLocked(domain string) bool
	Lock(domain string) error
}

type ThrottleService interface {
	CanCrawlURL(url string) bool
}

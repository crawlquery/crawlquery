package mem

import (
	"crawlquery/api/domain"
	"testing"

	"github.com/google/uuid"
)

func TestLock(t *testing.T) {
	t.Run("can lock domain", func(t *testing.T) {
		repo := NewRepository()
		key, err := repo.Lock("domain")
		if err != nil {
			t.Errorf("expected Lock to return nil, got %v", err)
		}

		if key == "" {
			t.Errorf("expected Lock to return key, got empty string")
		}

		if uuid.Validate(key) != nil {
			t.Errorf("expected Lock to return valid UUID, got %s", key)
		}
	})

	t.Run("cannot lock domain if already locked", func(t *testing.T) {
		repo := NewRepository()
		_, _ = repo.Lock("domain")
		_, err := repo.Lock("domain")
		if err != domain.ErrDomainLocked {
			t.Errorf("expected Lock to return ErrDomainLocked, got %v", err)
		}
	})

	t.Run("can unlock domain", func(t *testing.T) {
		repo := NewRepository()
		key, _ := repo.Lock("domain")
		err := repo.Unlock("domain", key)
		if err != nil {
			t.Errorf("expected Unlock to return nil, got %v", err)
		}
	})

	t.Run("cannot unlock domain if key is invalid", func(t *testing.T) {
		repo := NewRepository()
		_, _ = repo.Lock("domain")
		err := repo.Unlock("domain", "invalid")
		if err != domain.ErrInvalidLockKey {
			t.Errorf("expected Unlock to return ErrInvalidLockKey, got %v", err)
		}
	})
}

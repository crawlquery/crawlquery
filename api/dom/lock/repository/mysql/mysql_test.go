package mysql_test

import (
	"crawlquery/api/dom/lock/repository/mysql"
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"testing"

	"github.com/google/uuid"
)

func TestLock(t *testing.T) {
	t.Run("can lock domain", func(t *testing.T) {

		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		repo := mysql.NewRepository(db)

		defer db.Exec("DELETE FROM domain_locks WHERE domain = 'testlockdomain.com'")

		key, err := repo.Lock("testlockdomain.com")
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
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		repo := mysql.NewRepository(db)

		defer db.Exec("DELETE FROM domain_locks WHERE domain = 'cannotlockdomain.com'")

		_, _ = repo.Lock("cannotlockdomain.com")
		_, err := repo.Lock("cannotlockdomain.com")
		if err != domain.ErrDomainLocked {
			t.Errorf("expected Lock to return ErrDomainLocked, got %v", err)
		}
	})

	t.Run("can unlock domain", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		repo := mysql.NewRepository(db)

		defer db.Exec("DELETE FROM domain_locks WHERE domain = 'unlockdomain.com'")

		key, _ := repo.Lock("unlockdomain.com")
		err := repo.Unlock("unlockdomain.com", key)
		if err != nil {
			t.Errorf("expected Unlock to return nil, got %v", err)
		}
	})

	t.Run("cannot unlock domain if key is invalid", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		repo := mysql.NewRepository(db)

		defer db.Exec("DELETE FROM domain_locks WHERE domain = 'invalidkeydomain.com'")

		_, _ = repo.Lock("invalidkeydomain.com")
		err := repo.Unlock("invalidkeydomain.com", "invalid")
		if err != domain.ErrInvalidLockKey {
			t.Errorf("expected Unlock to return ErrInvalidLockKey, got %v", err)
		}
	})
}

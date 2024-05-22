package mysql_test

import (
	"crawlquery/api/crawl/restriction/repository/mysql"
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"database/sql"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	db := testutil.CreateTestMysqlDB()
	migration.Up(db)
	repo := mysql.NewRepository(db)

	t.Run("can get a restriction", func(t *testing.T) {
		restriction := &domain.CrawlRestriction{
			Domain: "can-get-a-restriction.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		}

		db.Exec("INSERT INTO crawl_restrictions (domain, until) VALUES (?, ?)", restriction.Domain, restriction.Until.Time)
		defer db.Exec("DELETE FROM crawl_restrictions WHERE domain = ?", restriction.Domain)

		got, err := repo.Get("can-get-a-restriction.com")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if got.Domain != restriction.Domain {
			t.Errorf("Expected %v, got %v", restriction, got)
		}

		if got.Until.Time.Sub(restriction.Until.Time) > time.Second {
			t.Errorf("Expected %v, got %v", restriction, got)
		}

		if got.Until.Valid != restriction.Until.Valid {
			t.Errorf("Expected %v, got %v", restriction, got)
		}
	})
}

func TestSet(t *testing.T) {
	db := testutil.CreateTestMysqlDB()
	migration.Up(db)
	repo := mysql.NewRepository(db)

	t.Run("can set a restriction", func(t *testing.T) {
		restriction := &domain.CrawlRestriction{
			Domain: "can-set-a-restriction.com",
			Until:  sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour)},
		}

		err := repo.Set(restriction)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		got, err := repo.Get("can-set-a-restriction.com")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if got.Domain != restriction.Domain {
			t.Errorf("Expected %v, got %v", restriction, got)
		}

		if got.Until.Time.Sub(restriction.Until.Time) > time.Second {
			t.Errorf("Expected %v, got %v", restriction, got)
		}

		if got.Until.Valid != restriction.Until.Valid {
			t.Errorf("Expected %v, got %v", restriction, got)
		}
	})
}

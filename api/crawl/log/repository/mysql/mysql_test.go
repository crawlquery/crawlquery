package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
	"time"

	crawlLogRepo "crawlquery/api/crawl/log/repository/mysql"
)

func TestSave(t *testing.T) {
	t.Run("should save a log", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()

		migration.Up(db)

		repo := crawlLogRepo.NewRepository(db)

		log := &domain.CrawlLog{
			ID:        domain.CrawlLogID(util.UUIDString()),
			PageID:    util.PageID("http://example.com"),
			Status:    domain.CrawlStatusCompleted,
			CreatedAt: time.Now(),
		}

		err := repo.Save(log)

		defer db.Exec("DELETE FROM crawl_logs WHERE id = ?", log.ID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		var savedLog domain.CrawlLog

		err = db.QueryRow("SELECT id, page_id, status FROM crawl_logs WHERE id = ?", log.ID).Scan(&savedLog.ID, &savedLog.PageID, &savedLog.Status)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if savedLog.PageID != log.PageID {
			t.Errorf("expected log pageID to be %v, got %v", log.PageID, savedLog.PageID)
		}

		if savedLog.Status != log.Status {
			t.Errorf("expected log status to be %v, got %v", log.Status, savedLog.Status)
		}
	})
}

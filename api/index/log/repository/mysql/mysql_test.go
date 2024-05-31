package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	indexLogRepo "crawlquery/api/index/log/repository/mysql"
)

func TestSave(t *testing.T) {
	t.Run("can save index log", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		repo := indexLogRepo.NewRepository(db)

		log := &domain.IndexLog{
			PageID:    "job1",
			Status:    domain.IndexStatusPending,
			Info:      "info",
			CreatedAt: time.Now(),
		}

		err := repo.Save(log)

		if err != nil {
			t.Errorf("Error creating index log: %v", err)
		}

		var checkLog domain.IndexLog

		err = db.QueryRow("SELECT id, page_id, status, info, created_at FROM index_logs WHERE page_id = ?", log.PageID).Scan(&checkLog.ID, &checkLog.PageID, &checkLog.Status, &checkLog.Info, &checkLog.CreatedAt)
		if err != nil {
			t.Errorf("Error getting index log: %v", err)
		}

		if checkLog.PageID != log.PageID {
			t.Errorf("Expected log ID to be %s, got %s", log.PageID, checkLog.PageID)
		}

		if checkLog.Status != log.Status {
			t.Errorf("Expected log status to be %s, got %s", log.Status, checkLog.Status)
		}

		if checkLog.Info != log.Info {
			t.Errorf("Expected log info to be %s, got %s", log.Info, checkLog.Info)
		}

		if checkLog.CreatedAt.UTC().Round(time.Minute) != log.CreatedAt.UTC().Round(time.Minute) {
			t.Errorf("Expected log created_at to be %s, got %s", log.CreatedAt, checkLog.CreatedAt)
		}
	})
}

func TestListByPageID(t *testing.T) {
	t.Run("can list index logs by pageID", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		repo := indexLogRepo.NewRepository(db)

		log1 := &domain.IndexLog{
			ID:        "log1",
			PageID:    "page1",
			Status:    domain.IndexStatusPending,
			Info:      "info",
			CreatedAt: time.Now(),
		}

		log2 := &domain.IndexLog{
			ID:        "log2",
			PageID:    "page2",
			Status:    domain.IndexStatusPending,
			Info:      "info",
			CreatedAt: time.Now(),
		}

		log3 := &domain.IndexLog{
			ID:        "log3",
			PageID:    "page3",
			Status:    domain.IndexStatusCompleted,
			Info:      "info",
			CreatedAt: time.Now(),
		}

		for _, log := range []*domain.IndexLog{log1, log2, log3} {
			_, err := db.Exec("INSERT INTO index_logs (id, page_id, status, info, created_at) VALUES (?, ?, ?, ?, ?)", log.ID, log.PageID, log.Status, log.Info, log.CreatedAt)
			if err != nil {
				t.Errorf("Error creating index log: %v", err)
			}
		}

		logs, err := repo.ListByPageID(log1.PageID)

		if err != nil {
			t.Errorf("Error listing index logs: %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("Expected 1 log, got %v", len(logs))
		}

		if logs[0].PageID != log1.PageID {
			t.Errorf("Expected log pageID to be %s, got %s", log1.PageID, logs[0].PageID)
		}

		if logs[0].Status != log1.Status {
			t.Errorf("Expected log status to be %s, got %s", log1.Status, logs[0].Status)
		}

		if logs[0].Info != log1.Info {
			t.Errorf("Expected log info to be %s, got %s", log1.Info, logs[0].Info)
		}
	})
}

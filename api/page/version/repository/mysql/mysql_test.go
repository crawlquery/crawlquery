package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
	"time"

	pageVersionRepo "crawlquery/api/page/version/repository/mysql"
)

func TestGet(t *testing.T) {
	t.Run("returns a page version", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := pageVersionRepo.NewRepository(db)
		pageVersion := &domain.PageVersion{
			ID:          domain.PageVersionID(util.UUIDString()),
			PageID:      "page1",
			ContentHash: "hash",
			CreatedAt:   time.Now(),
		}

		_, err := db.Exec("INSERT INTO page_versions (id, page_id, content_hash, created_at) VALUES (?, ?, ?, ?)", pageVersion.ID, pageVersion.PageID, pageVersion.ContentHash, pageVersion.CreatedAt)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defer db.Exec("DELETE FROM page_versions WHERE id = ?", pageVersion.ID)

		got, err := repo.Get(pageVersion.ID)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if got.ID != pageVersion.ID {
			t.Fatalf("Expected %s, got %s", pageVersion.ID, got.ID)
		}

		if got.PageID != pageVersion.PageID {
			t.Fatalf("Expected %s, got %s", pageVersion.PageID, got.PageID)
		}

		if got.ContentHash != pageVersion.ContentHash {
			t.Fatalf("Expected %s, got %s", pageVersion.ContentHash, got.ContentHash)
		}
	})
}

func TestListByPageID(t *testing.T) {
	t.Run("returns page versions", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := pageVersionRepo.NewRepository(db)
		pageID := domain.PageID("page1")

		pageVersion1 := &domain.PageVersion{
			ID:          domain.PageVersionID(util.UUIDString()),
			PageID:      pageID,
			ContentHash: "hash1",
			CreatedAt:   time.Now(),
		}

		pageVersion2 := &domain.PageVersion{
			ID:          domain.PageVersionID(util.UUIDString()),
			PageID:      pageID,
			ContentHash: "hash2",
			CreatedAt:   time.Now(),
		}

		_, err := db.Exec("INSERT INTO page_versions (id, page_id, content_hash, created_at) VALUES (?, ?, ?, ?)", pageVersion1.ID, pageVersion1.PageID, pageVersion1.ContentHash, pageVersion1.CreatedAt)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defer db.Exec("DELETE FROM page_versions WHERE id = ?", pageVersion1.ID)

		_, err = db.Exec("INSERT INTO page_versions (id, page_id, content_hash, created_at) VALUES (?, ?, ?, ?)", pageVersion2.ID, pageVersion2.PageID, pageVersion2.ContentHash, pageVersion2.CreatedAt)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defer db.Exec("DELETE FROM page_versions WHERE id = ?", pageVersion2.ID)

		got, err := repo.ListByPageID(pageID)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("Expected 2, got %d", len(got))
		}

		if got[0].ID != pageVersion1.ID {
			t.Fatalf("Expected %s, got %s", pageVersion1.ID, got[0].ID)
		}

		if got[1].ID != pageVersion2.ID {
			t.Fatalf("Expected %s, got %s", pageVersion2.ID, got[1].ID)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("creates a page version", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := pageVersionRepo.NewRepository(db)
		pageVersion := &domain.PageVersion{
			ID:          domain.PageVersionID(util.UUIDString()),
			PageID:      "page1",
			ContentHash: "hash",
			CreatedAt:   time.Now(),
		}

		err := repo.Create(pageVersion)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defer db.Exec("DELETE FROM page_versions WHERE id = ?", pageVersion.ID)

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM page_versions WHERE id = ?", pageVersion.ID).Scan(&count)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if count != 1 {
			t.Fatalf("Expected 1, got %d", count)
		}

		got, err := repo.Get(pageVersion.ID)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if got.ID != pageVersion.ID {
			t.Fatalf("Expected %s, got %s", pageVersion.ID, got.ID)
		}

		if got.PageID != pageVersion.PageID {
			t.Fatalf("Expected %s, got %s", pageVersion.PageID, got.PageID)
		}

		if got.ContentHash != pageVersion.ContentHash {
			t.Fatalf("Expected %s, got %s", pageVersion.ContentHash, got.ContentHash)
		}

		if got.CreatedAt.Round(time.Second) != pageVersion.CreatedAt.UTC().Round(time.Second) {
			t.Fatalf("Expected %v, got %v", pageVersion.CreatedAt, got.CreatedAt)
		}

	})
}

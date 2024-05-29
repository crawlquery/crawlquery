package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/api/page/repository/mysql"
	"crawlquery/pkg/testutil"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Run("should return a page", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		err := migration.Up(db)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		repo := mysql.NewRepository(db)

		page := &domain.Page{
			ID:        "123",
			ShardID:   1,
			URL:       "http://example.com",
			CreatedAt: time.Now(),
		}

		_, err = db.Exec("INSERT INTO pages (id, url, shard_id, created_at) VALUES (?, ?, ?, ?)", page.ID, page.URL, page.ShardID, page.CreatedAt)

		defer db.Exec("DELETE FROM pages WHERE id = ?", page.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		res, err := repo.Get("123")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if res.ID != "123" {
			t.Errorf("expected page ID to be 123, got %s", res.ID)
		}

		if res.ShardID != 1 {
			t.Errorf("expected page ShardID to be 1, got %d", res.ShardID)
		}

		if res.URL != "http://example.com" {
			t.Errorf("expected page URL to be http://example.com, got %s", res.URL)
		}

		if res.CreatedAt.UTC().Round(time.Second) != page.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("expected page CreatedAt to be %v, got %v", page.CreatedAt, res.CreatedAt)
		}
	})

	t.Run("should return ErrPageNotFound if page not found", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		_, err := repo.Get("123")

		if err != domain.ErrPageNotFound {
			t.Errorf("expected error to be ErrPageNotFound, got %v", err)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("should save a page", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		page := &domain.Page{
			ID:        "123",
			URL:       "http://example.com",
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		err := repo.Create(page)

		defer db.Exec("DELETE FROM pages WHERE id = ?", page.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var id domain.PageID
		var url domain.URL
		var shardID domain.ShardID
		var createdAt time.Time

		err = db.QueryRow("SELECT id, url, shard_id, created_at FROM pages WHERE id = ?", page.ID).Scan(&id, &url, &shardID, &createdAt)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if id != "123" {
			t.Errorf("expected page ID to be 123, got %s", id)
		}

		if url != "http://example.com" {
			t.Errorf("expected page URL to be http://example.com, got %s", url)
		}

		if shardID != 1 {
			t.Errorf("expected page ShardID to be 1, got %d", shardID)
		}

		if createdAt.UTC().Round(time.Second) != page.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("expected page CreatedAt to be %v, got %v", page.CreatedAt, createdAt)
		}
	})
}

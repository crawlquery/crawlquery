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
		migration.Up(db)

		repo := mysql.NewRepository(db)

		page := &domain.Page{
			ID:        "123",
			ShardID:   1,
			Hash:      "abc123",
			CreatedAt: time.Now(),
		}

		_, err := db.Exec("INSERT INTO pages (id, shard_id, hash, created_at) VALUES (?, ?, ?, ?)", page.ID, page.ShardID, page.Hash, page.CreatedAt)

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

		if res.Hash != "abc123" {
			t.Errorf("expected page Hash to be abc123, got %s", res.Hash)
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
			ShardID:   1,
			Hash:      "abc123",
			CreatedAt: time.Now(),
		}

		err := repo.Create(page)

		defer db.Exec("DELETE FROM pages WHERE id = ?", page.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var id string
		var shardID int
		var hash string
		var createdAt time.Time

		err = db.QueryRow("SELECT id, shard_id, hash, created_at FROM pages WHERE id = ?", page.ID).Scan(&id, &shardID, &hash, &createdAt)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if id != "123" {
			t.Errorf("expected page ID to be 123, got %s", id)
		}

		if shardID != 1 {
			t.Errorf("expected page ShardID to be 1, got %d", shardID)
		}

		if hash != "abc123" {
			t.Errorf("expected page Hash to be abc123, got %s", hash)
		}

		if createdAt.UTC().Round(time.Second) != page.CreatedAt.UTC().Round(time.Second) {
			t.Errorf("expected page CreatedAt to be %v, got %v", page.CreatedAt, createdAt)
		}
	})
}

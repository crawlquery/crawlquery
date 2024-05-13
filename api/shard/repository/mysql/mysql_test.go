package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/api/shard/repository/mysql"
	"crawlquery/pkg/testutil"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create a shard", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		shard := &domain.Shard{
			ID:        3,
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM shards WHERE id = ?", shard.ID)

		err := repo.Create(shard)

		// Assert
		if err != nil {
			t.Errorf("Error creating shards: %v", err)
		}

		res, err := db.Query("SELECT id, created_at FROM shards WHERE id = ?", shard.ID)

		if err != nil {
			t.Errorf("Error querying for shard: %v", err)
		}

		var id uint
		var createdAt time.Time

		for res.Next() {
			err = res.Scan(&id, &createdAt)
			if err != nil {
				t.Errorf("Error scanning account: %v", err)
			}
		}

		if id != shard.ID {
			t.Errorf("Expected ID to be %d, got %d", shard.ID, id)
		}

		if createdAt.Sub(shard.CreatedAt) > time.Second || shard.CreatedAt.Sub(createdAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", shard.CreatedAt, createdAt)
		}
	})
}

func TestList(t *testing.T) {
	t.Run("can list shards", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		shard := &domain.Shard{
			ID:        3,
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM shards WHERE id = ?", shard.ID)

		err := repo.Create(shard)

		if err != nil {
			t.Fatalf("Error creating shard: %v", err)
		}

		// Act
		list, err := repo.List()

		// Assert
		if err != nil {
			t.Fatalf("Error listing shards: %v", err)
		}

		if len(list) != 1 {
			t.Fatalf("Expected 1 shard, got %d", len(list))
		}

		if list[0].ID != shard.ID {
			t.Errorf("Expected ID to be %d, got %d", shard.ID, list[0].ID)
		}

		if list[0].CreatedAt.Sub(shard.CreatedAt) > time.Second || shard.CreatedAt.Sub(list[0].CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", shard.CreatedAt, list[0].CreatedAt)
		}
	})
}

func TestCount(t *testing.T) {
	t.Run("can count shards", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		shard := &domain.Shard{
			ID:        3,
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM shards WHERE id = ?", shard.ID)

		err := repo.Create(shard)

		if err != nil {
			t.Fatalf("Error creating shard: %v", err)
		}

		// Act
		count, err := repo.Count()

		// Assert
		if err != nil {
			t.Fatalf("Error counting shards: %v", err)
		}

		if count != 1 {
			t.Fatalf("Expected 1 shard, got %d", count)
		}
	})
}

package mem

import (
	"crawlquery/api/domain"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create a shard", func(t *testing.T) {
		// Arrange
		repo := NewRepository()

		shard := &domain.Shard{
			ID:        domain.ShardID(3),
			CreatedAt: time.Now(),
		}

		err := repo.Create(shard)

		if err != nil {
			t.Fatalf("Error creating shard: %v", err)
		}

		check, ok := repo.shards[shard.ID]

		if !ok {
			t.Fatalf("Expected shard to be in repository")
		}

		if check.CreatedAt.Sub(shard.CreatedAt) > time.Second || shard.CreatedAt.Sub(check.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", shard.CreatedAt, check.CreatedAt)
		}
	})

}

func TestGet(t *testing.T) {
	t.Run("can get a shard", func(t *testing.T) {
		// Arrange
		repo := NewRepository()

		shard := &domain.Shard{
			ID:        domain.ShardID(3),
			CreatedAt: time.Now(),
		}

		repo.shards[shard.ID] = shard

		// Act
		check, err := repo.Get(shard.ID)

		// Assert
		if err != nil {
			t.Fatalf("Error getting shard: %v", err)
		}

		if check.ID != shard.ID {
			t.Errorf("Expected ID to be %d, got %d", shard.ID, check.ID)
		}

		if check.CreatedAt.Sub(shard.CreatedAt) > time.Second || shard.CreatedAt.Sub(check.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", shard.CreatedAt, check.CreatedAt)
		}
	})
}

func TestList(t *testing.T) {
	t.Run("can list shards", func(t *testing.T) {
		// Arrange
		repo := NewRepository()

		shard := &domain.Shard{
			ID:        domain.ShardID(3),
			CreatedAt: time.Now(),
		}

		repo.shards[shard.ID] = shard

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
		repo := NewRepository()

		shard := &domain.Shard{
			ID:        3,
			CreatedAt: time.Now(),
		}

		repo.shards[shard.ID] = shard

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

package service_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/shard/repository/mem"
	"crawlquery/api/shard/service"
	"crawlquery/pkg/testutil"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create a shard", func(t *testing.T) {
		repo := mem.NewRepository()
		service := service.NewService(repo, testutil.NewTestLogger())

		shard := &domain.Shard{
			ID: 3,
		}

		err := service.Create(shard)

		if err != nil {
			t.Fatalf("Error creating shard: %v", err)
		}

		check, err := repo.Get(shard.ID)

		if err != nil {
			t.Fatalf("Error getting shard: %v", err)
		}

		if check.ID != shard.ID {
			t.Errorf("Expected ID to be %d, got %d", shard.ID, check.ID)
		}
	})
}

func TestGetURLShardID(t *testing.T) {
	// Define test cases with URLs and the expected shard for a given number of shards
	tests := []struct {
		url        string
		numShards  int
		expectedID uint
	}{
		{"https://www.amazon.com", 5000, 4786},
		{"https://www.google.com", 10, 5},
		{"https://www.example.com", 10, 8},
		{"https://openai.com", 10, 5},
		{"https://www.randomsite.org", 10, 9},
		{"https://www.differentnumber.com", 5, 3},
		{"https://www.anotherone.com", 5, 1},
		{"https://www.lastone.com", 1, 0},
	}

	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			repo := mem.NewRepository()
			service := service.NewService(repo, testutil.NewTestLogger())

			for i := 0; i < tc.numShards; i++ {
				repo.Create(&domain.Shard{ID: uint(i)})
			}

			shardID, err := service.GetURLShardID(tc.url)

			if err != nil {
				t.Errorf("Error getting shard ID: %v", err)
			}

			if shardID != tc.expectedID {
				t.Errorf("getShardID(%q, %d) = %d; want %d", tc.url, tc.numShards, shardID, tc.expectedID)
			}
		})
	}
}

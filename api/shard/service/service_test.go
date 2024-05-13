package service_test

import (
	"crawlquery/api/domain"
	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"
	"crawlquery/api/shard/repository/mem"
	"crawlquery/api/shard/service"
	"crawlquery/pkg/testutil"
	"testing"
)

func TestGetURLShardID(t *testing.T) {
	// Define test cases with URLs and the expected shard for a given number of shards
	tests := []struct {
		url        string
		numShards  int
		expectedID int
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
			service := service.NewService(repo, nil, testutil.NewTestLogger())

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

func TestGetShardWithLeastNodes(t *testing.T) {
	t.Run("can get shard with least nodes", func(t *testing.T) {
		repo := mem.NewRepository()
		nodeRepo := nodeRepo.NewRepository()
		nodeService := nodeService.NewService(nodeRepo, nil, testutil.NewTestLogger())

		nodes := []*domain.Node{
			{ID: "1", ShardID: 1},
			{ID: "2", ShardID: 1},
			{ID: "3", ShardID: 2},
			{ID: "4", ShardID: 2},
			{ID: "5", ShardID: 3},
		}

		for _, n := range nodes {
			nodeRepo.Create(n)
		}

		service := service.NewService(repo, nodeService, testutil.NewTestLogger())

		shard, err := service.GetShardWithLeastNodes()

		if err != nil {
			t.Fatalf("Error getting shard with least nodes: %v", err)
		}

		if shard.ID != 3 {
			t.Errorf("Expected shard ID to be 3, got %d", shard.ID)
		}
	})
}

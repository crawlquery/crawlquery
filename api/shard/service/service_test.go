package service_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/shard/repository/mem"
	"crawlquery/api/shard/service"
	"crawlquery/pkg/testutil"
	"testing"
)

func TestGetURLShardID(t *testing.T) {
	// Define test cases with URLs and the expected shard for a given number of shards
	tests := []struct {
		url        domain.URL
		numShards  int
		expectedID domain.ShardID
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
		t.Run(string(tc.url), func(t *testing.T) {
			repo := mem.NewRepository()
			service := service.NewService(
				service.WithRepo(repo),
				service.WithLogger(testutil.NewTestLogger()),
			)

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

package service_test

import (
	"crawlquery/api/shard/service"
	"testing"
)

func TestGetShardID(t *testing.T) {
	// Define test cases with URLs and the expected shard for a given number of shards
	tests := []struct {
		url        string
		numShards  int
		expectedID int
	}{
		{"https://www.google.com", 10, service.GetShardID("https://www.google.com", 10)},
		{"https://www.example.com", 10, service.GetShardID("https://www.example.com", 10)},
		{"https://openai.com", 10, service.GetShardID("https://openai.com", 10)},
		{"https://www.randomsite.org", 10, service.GetShardID("https://www.randomsite.org", 10)},
		{"https://www.differentnumber.com", 5, service.GetShardID("https://www.differentnumber.com", 5)},
		{"https://www.anotherone.com", 5, service.GetShardID("https://www.anotherone.com", 5)},
	}

	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			shardID := service.GetShardID(tc.url, tc.numShards)
			if shardID != tc.expectedID {
				t.Errorf("getShardID(%q, %d) = %d; want %d", tc.url, tc.numShards, shardID, tc.expectedID)
			}
		})
	}
}

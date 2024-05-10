package shard

import "hash/fnv"

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func GetShardID(url string, numShards int) int {
	return int(hash(url) % uint32(numShards))
}

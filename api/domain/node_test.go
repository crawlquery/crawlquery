package domain_test

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"strings"
	"testing"
	"time"
)

func TestNodeValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		err := node.Validate()

		if err != nil {
			t.Errorf("Expected node to be valid, got error: %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		node := &domain.Node{
			ID:        "aaa",
			AccountID: util.UUID(),
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		err := node.Validate()

		if err == nil {
			t.Errorf("Expected node to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Node.ID") {
			t.Errorf("Expected error to contain 'Node.ID', got %v", err)
		}
	})

	t.Run("invalid hostname", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "!!",
		}

		err := node.Validate()

		if err == nil {
			t.Errorf("Expected node to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Node.Hostname") {
			t.Errorf("Expected error to contain 'Node.Hostname', got %v", err)
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "testnode",
			Port:      200000,
		}

		err := node.Validate()

		if err == nil {
			t.Errorf("Expected node to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Node.Port") {
			t.Errorf("Expected error to contain 'Node.Port', got %v", err)
		}
	})

	t.Run("invalid shard id", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   2000000,
		}

		err := node.Validate()

		if err == nil {
			t.Errorf("Expected node to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Node.ShardID") {
			t.Errorf("Expected error to contain 'Node.ShardID', got %v", err)
		}
	})
}

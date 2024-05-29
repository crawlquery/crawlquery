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
			ID:        util.UUIDString(),
			Key:       util.UUIDString(),
			AccountID: util.UUIDString(),
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

	t.Run("key is required", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: util.UUIDString(),
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		err := node.Validate()

		if err == nil {
			t.Errorf("Expected node to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Node.Key") {
			t.Errorf("Expected error to contain 'Node.Key', got %v", err)
		}
	})

	t.Run("key is uuid", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUIDString(),
			Key:       "fails",
			AccountID: util.UUIDString(),
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		err := node.Validate()

		if err == nil {
			t.Errorf("Expected node to be invalid, got nil")
		}

		if !strings.Contains(err.Error(), "Node.Key") {
			t.Errorf("Expected error to contain 'Node.Key', got %v", err)
		}
	})

	t.Run("key is valid", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUIDString(),
			Key:       util.UUIDString(),
			AccountID: util.UUIDString(),
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
			AccountID: util.UUIDString(),
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
			ID:        util.UUIDString(),
			AccountID: util.UUIDString(),
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
			ID:        util.UUIDString(),
			AccountID: util.UUIDString(),
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

}

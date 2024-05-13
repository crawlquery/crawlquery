package dto_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreateNodeRequestToJSON(t *testing.T) {
	t.Run("should return correct JSON", func(t *testing.T) {
		req := &dto.CreateNodeRequest{
			Hostname: "localhost",
			Port:     8080,
		}

		b, err := req.ToJSON()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := `{"hostname":"localhost","port":8080}`
		if string(b) != expected {
			t.Errorf("expected: %s, got: %s", expected, string(b))
		}
	})
}

func TestNewCreateNodeResponse(t *testing.T) {
	t.Run("should return correct CreateNodeResponse from Node", func(t *testing.T) {
		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: util.UUID(),
			Hostname:  "localhost",
			Port:      8080,
			ShardID:   1,
		}

		resp := dto.NewCreateNodeResponse(node)

		if resp.Node.ID != node.ID {
			t.Errorf("expected: %s, got: %s", node.ID, resp.Node.ID)
		}

		if resp.Node.AccountID != node.AccountID {
			t.Errorf("expected: %s, got: %s", node.AccountID, resp.Node.AccountID)
		}

		if resp.Node.Hostname != node.Hostname {
			t.Errorf("expected: %s, got: %s", node.Hostname, resp.Node.Hostname)
		}

		if resp.Node.Port != node.Port {
			t.Errorf("expected: %d, got: %d", node.Port, resp.Node.Port)
		}

		if resp.Node.ShardID != node.ShardID {
			t.Errorf("expected: %d, got: %d", node.ShardID, resp.Node.ShardID)
		}

	})
}

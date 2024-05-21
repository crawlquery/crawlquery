package api_test

import (
	"crawlquery/api/dto"
	"crawlquery/pkg/client/api"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestAuthenticateNode(t *testing.T) {
	t.Run("should return a node if the key is correct", func(t *testing.T) {
		mockRes := &dto.AuthenticateNodeResponse{
			Node: &dto.Node{
				ID:        "123",
				Key:       "123",
				AccountID: "123",
				Hostname:  "localhost",
				Port:      8080,
				ShardID:   1,
				CreatedAt: time.Now(),
			},
		}

		defer gock.Off()

		gock.New("http://localhost:8080").
			Post("auth/node").
			JSON(`{"key":"123"}`).
			Reply(200).
			JSON(mockRes)

		client := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		resNode, err := client.AuthenticateNode("123")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if resNode.ID != mockRes.Node.ID {
			t.Errorf("expected node ID: %s, got: %s", mockRes.Node.ID, resNode.ID)
		}

		if resNode.Key != mockRes.Node.Key {
			t.Errorf("expected node key: %s, got: %s", mockRes.Node.Key, resNode.Key)
		}

		if resNode.AccountID != mockRes.Node.AccountID {
			t.Errorf("expected node account ID: %s, got: %s", mockRes.Node.AccountID, resNode.AccountID)
		}

		if resNode.Hostname != mockRes.Node.Hostname {
			t.Errorf("expected node hostname: %s, got: %s", mockRes.Node.Hostname, resNode.Hostname)
		}

		if resNode.Port != mockRes.Node.Port {
			t.Errorf("expected node port: %d, got: %d", mockRes.Node.Port, resNode.Port)
		}

		if resNode.ShardID != mockRes.Node.ShardID {
			t.Errorf("expected node shard ID: %d, got: %d", mockRes.Node.ShardID, resNode.ShardID)
		}

		if resNode.CreatedAt.IsZero() {
			t.Errorf("expected node created at to be set, got zero value")
		}

	})

	t.Run("should return an error if the key is incorrect", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://localhost:8080").
			Post("auth/node").
			JSON(`{"key":"123"}`).
			Reply(401)

		client := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		_, err := client.AuthenticateNode("123")

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != "could not authenticate node" {
			t.Errorf("expected error: could not authenticate node, got: %v", err)
		}
	})
}

func TestListNodesByShardID(t *testing.T) {
	t.Run("should return nodes by shard ID", func(t *testing.T) {
		mockRes := &dto.ListNodesResponse{
			Nodes: []*dto.Node{
				{
					ID:        "123",
					Key:       "123",
					AccountID: "123",
					Hostname:  "localhost",
					Port:      8080,
					ShardID:   1,
					CreatedAt: time.Now(),
				},
				{
					ID:        "456",
					Key:       "456",
					AccountID: "456",
					Hostname:  "localhost",
					Port:      8080,
					ShardID:   1,
					CreatedAt: time.Now(),
				},
			},
		}

		defer gock.Off()

		gock.New("http://localhost:8080").
			Get("shards/1/nodes").
			Reply(200).
			JSON(mockRes)

		client := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		resNodes, err := client.ListNodesByShardID(1)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(resNodes) != len(mockRes.Nodes) {
			t.Errorf("expected %d nodes, got %d", len(mockRes.Nodes), len(resNodes))
		}

		for i, resNode := range resNodes {
			mockNode := mockRes.Nodes[i]

			if resNode.ID != mockNode.ID {
				t.Errorf("expected node ID: %s, got: %s", mockNode.ID, resNode.ID)
			}

			if resNode.Hostname != mockNode.Hostname {
				t.Errorf("expected node hostname: %s, got: %s", mockNode.Hostname, resNode.Hostname)
			}

			if resNode.Port != mockNode.Port {
				t.Errorf("expected node port: %d, got: %d", mockNode.Port, resNode.Port)
			}

			if resNode.ShardID != mockNode.ShardID {
				t.Errorf("expected node shard ID: %d, got: %d", mockNode.ShardID, resNode.ShardID)
			}
		}
	})
}

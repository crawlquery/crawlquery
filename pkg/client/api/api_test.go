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
		node := &dto.Node{
			ID:        "123",
			Key:       "123",
			AccountID: "123",
			Hostname:  "localhost",
			Port:      8080,
			ShardID:   1,
			CreatedAt: time.Now(),
		}

		defer gock.Off()

		gock.New("http://localhost:8080").
			Post("auth/node").
			JSON(`{"key":"123"}`).
			Reply(200).
			JSON(node)

		client := api.NewClient("http://localhost:8080", testutil.NewTestLogger())

		resNode, err := client.AuthenticateNode("123")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if resNode.ID != node.ID {
			t.Errorf("expected: %s, got: %s", node.ID, resNode.ID)
		}

		if resNode.Key != node.Key {
			t.Errorf("expected: %s, got: %s", node.Key, resNode.Key)
		}

		if resNode.AccountID != node.AccountID {
			t.Errorf("expected: %s, got: %s", node.AccountID, resNode.AccountID)
		}

		if resNode.Hostname != node.Hostname {
			t.Errorf("expected: %s, got: %s", node.Hostname, resNode.Hostname)
		}

		if resNode.Port != node.Port {
			t.Errorf("expected: %d, got: %d", node.Port, resNode.Port)
		}

		if resNode.ShardID != node.ShardID {
			t.Errorf("expected: %d, got: %d", node.ShardID, resNode.ShardID)
		}

		if resNode.CreatedAt.IsZero() {
			t.Errorf("expected CreatedAt to be set, got zero value")
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

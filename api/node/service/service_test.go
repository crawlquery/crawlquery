package service_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/factory"
	"crawlquery/pkg/util"
	"errors"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create a node", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, repo := factory.NodeService(accSvc)

		node, err := svc.Create(
			accountID,
			"testnode",
			8080,
		)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		if node.AccountID != accountID {
			t.Errorf("Expected AccountID to be %s, got %s", accountID, node.AccountID)
		}

		if node.Hostname != "testnode" {
			t.Errorf("Expected Hostname to be 'testnode', got %s", node.Hostname)
		}

		if node.Port != 8080 {
			t.Errorf("Expected Port to be 8080, got %d", node.Port)
		}

		list, err := repo.List()

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(list) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(list))
		}

		if list[0].Hostname != "testnode" {
			t.Errorf("Expected Hostname to be 'testnode', got %s", list[0].Hostname)
		}

		if list[0].Port != 8080 {
			t.Errorf("Expected Port to be 8080, got %d", list[0].Port)
		}
	})

	t.Run("can't create a node that already exists", func(t *testing.T) {
		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, _ := factory.NodeService(accSvc)

		_, err := svc.Create(
			accountID,
			"hostname",
			8080,
		)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		_, err = svc.Create(
			accountID,
			"hostname",
			8080,
		)

		if err != domain.ErrNodeAlreadyExists {
			t.Errorf("Expected error creating node with same hostname")
		}
	})

	t.Run("can't create a node with AccountID that doesn't exist", func(t *testing.T) {
		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(nil)

		svc, _ := factory.NodeService(accSvc)

		_, err := svc.Create(
			accountID,
			"hostname",
			8080,
		)

		if err != domain.ErrInvalidAccountID {
			t.Errorf("Expected error creating node with invalid AccountID, got %v", err)
		}
	})

	t.Run("can't create a node with invalid hostname", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, _ := factory.NodeService(accSvc)

		_, err := svc.Create(
			accountID,
			"!!",
			8080,
		)

		if err == nil {
			t.Errorf("Expected error creating node with invalid hostname")
		}
	})

	t.Run("handles error creating node", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, repo := factory.NodeService(accSvc)

		repo.ForceCreateError(errors.New("db locked"))

		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: accountID,
			Hostname:  "testnode",
		}

		repo.Create(node)

		_, err := svc.Create(
			accountID,
			"testnode",
			8080,
		)

		if err == nil {
			t.Errorf("Expected error creating node with same hostname")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("can list nodes", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, _ := factory.NodeService(accSvc)

		node, _ := svc.Create(
			accountID,
			"testnode",
			8080,
		)

		node2, _ := svc.Create(
			accountID,
			"testnode2",
			8081,
		)

		list, err := svc.List()

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(list) != 2 {
			t.Fatalf("Expected 2 nodes, got %d", len(list))
		}

		for _, n := range list {
			if n.Hostname != node.Hostname && n.Hostname != node2.Hostname {
				t.Errorf("Expected node to be one of %s or %s, got %s", node.Hostname, node2.Hostname, n.Hostname)
			}

			if n.Port != node.Port && n.Port != node2.Port {
				t.Errorf("Expected port to be one of %d or %d, got %d", node.Port, node2.Port, n.Port)
			}
		}
	})
}

func TestRandomizedList(t *testing.T) {
	t.Run("can list nodes in random order", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, repo := factory.NodeService(accSvc)

		node := &domain.Node{
			ID:        util.UUID(),
			AccountID: accountID,
			Hostname:  "testnode",
		}

		node2 := &domain.Node{
			ID:        util.UUID(),
			AccountID: accountID,
			Hostname:  "testnode2",
		}

		node3 := &domain.Node{
			ID:        util.UUID(),
			AccountID: accountID,
			Hostname:  "testnode3",
		}

		repo.Create(node)
		repo.Create(node2)
		repo.Create(node3)

		list, err := svc.RandomizedList()

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(list) != 3 {
			t.Fatalf("Expected 3 nodes, got %d", len(list))
		}

		var firstSeenCount int

		for i := 100; i > 0; i-- {
			randList, _ := svc.RandomizedList()

			if list[0].ID == randList[0].ID {
				firstSeenCount++
			}
		}

		if firstSeenCount > 90 {
			t.Errorf("Expected first node to be in a different position at least once")
		}
	})

	t.Run("handles error listing nodes", func(t *testing.T) {

		accountID := util.UUID()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, repo := factory.NodeService(accSvc)

		repo.ForceListError(errors.New("db locked"))

		_, err := svc.RandomizedList()

		if err == nil {
			t.Errorf("Expected error listing nodes")
		}
	})
}

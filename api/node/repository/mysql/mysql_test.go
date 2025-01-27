package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/api/node/repository/mysql"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create a node", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		err := migration.Up(db)
		if err != nil {
			t.Fatalf("Error running migrations: %v", err)
		}

		repo := mysql.NewRepository(db)
		node := &domain.Node{
			ID:        util.UUIDString(),
			Key:       util.UUIDString(),
			AccountID: util.UUIDString(),
			Hostname:  "testnode",
		}
		defer db.Exec("DELETE FROM nodes WHERE id = ?", node.ID)

		err = repo.Create(node)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		var check domain.Node
		err = db.QueryRow("SELECT id, `key`, account_id, hostname, port, shard_id, created_at FROM nodes WHERE id = ?", node.ID).Scan(
			&check.ID,
			&check.Key,
			&check.AccountID,
			&check.Hostname,
			&check.Port,
			&check.ShardID,
			&check.CreatedAt,
		)

		if err != nil {
			t.Fatalf("Error getting node: %v", err)
		}

		if check.ID != node.ID {
			t.Errorf("Expected ID to be %s, got %s", node.ID, check.ID)
		}

		if check.AccountID != node.AccountID {
			t.Errorf("Expected AccountID to be %s, got %s", node.AccountID, check.AccountID)
		}

		if check.Key != node.Key {
			t.Errorf("Expected Key to be %s, got %s", node.Key, check.Key)
		}

		if check.Hostname != node.Hostname {
			t.Errorf("Expected Name to be %s, got %s", node.Hostname, check.Hostname)
		}

		if check.Port != node.Port {
			t.Errorf("Expected Port to be %d, got %d", node.Port, check.Port)
		}

		if check.ShardID != node.ShardID {
			t.Errorf("Expected ShardID to be %d, got %d", node.ShardID, check.ShardID)
		}

		if check.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set, got zero value")
		}
	})

	t.Run("can't create a node with the same hostname and port", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)
		node := &domain.Node{
			ID:       util.UUIDString(),
			Hostname: "testnode",
			Port:     8080,
		}

		defer db.Exec("DELETE FROM nodes WHERE id = ?", node.ID)

		err := repo.Create(node)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		err = repo.Create(node)

		if err == nil {
			t.Fatalf("Expected error creating node with same hostname, got nil")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("can list nodes", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)
		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: util.UUIDString(),
			Key:       util.UUIDString(),
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   1,
		}

		node2 := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: util.UUIDString(),
			Key:       util.UUIDString(),
			Hostname:  "testnode2",
			Port:      8081,
			ShardID:   2,
		}

		defer db.Exec("DELETE FROM nodes WHERE id = ?", node.ID)
		defer db.Exec("DELETE FROM nodes WHERE id = ?", node2.ID)

		err := repo.Create(node)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		err = repo.Create(node2)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		nodes, err := repo.List()

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		for _, n := range nodes {
			if n.Hostname != node.Hostname && n.Hostname != node2.Hostname {
				t.Errorf("Expected node to be one of %s or %s, got %s", node.Hostname, node2.Hostname, n.Hostname)
			}

			if n.Port != node.Port && n.Port != node2.Port {
				t.Errorf("Expected port to be one of %d or %d, got %d", node.Port, node2.Port, n.Port)
			}
		}
	})
}

func TestListByAccounID(t *testing.T) {
	t.Run("can list nodes by account id", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)
		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: "account1",
			Hostname:  "testnode",
			Key:       util.UUIDString(),
			Port:      8080,
			ShardID:   1,
		}

		node2 := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: "account2",
			Hostname:  "testnode2",
			Key:       util.UUIDString(),
			Port:      8081,
			ShardID:   2,
		}

		defer db.Exec("DELETE FROM nodes WHERE id = ?", node.ID)
		defer db.Exec("DELETE FROM nodes WHERE id = ?", node2.ID)

		err := repo.Create(node)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		err = repo.Create(node2)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		nodes, err := repo.ListByAccountID("account1")

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(nodes) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(nodes))
		}

		if nodes[0].Hostname != node.Hostname {
			t.Errorf("Expected node to be %s, got %s", node.Hostname, nodes[0].Hostname)
		}

		if nodes[0].Port != node.Port {
			t.Errorf("Expected port to be %d, got %d", node.Port, nodes[0].Port)
		}
	})
}

func TestGetNodeByKey(t *testing.T) {
	t.Run("can get node by key", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)
		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: "account1",
			Hostname:  "testnode",
			Port:      8080,
			ShardID:   1,
		}

		defer db.Exec("DELETE FROM nodes WHERE id = ?", node.ID)

		err := repo.Create(node)

		if err != nil {
			t.Fatalf("Error creating node: %v", err)
		}

		check, err := repo.GetNodeByKey(node.Key)

		if err != nil {
			t.Fatalf("Error getting node: %v", err)
		}

		if check.ID != node.ID {
			t.Errorf("Expected ID to be %s, got %s", node.ID, check.ID)
		}

		if check.AccountID != node.AccountID {
			t.Errorf("Expected AccountID to be %s, got %s", node.AccountID, check.AccountID)
		}

		if check.Key != node.Key {
			t.Errorf("Expected Key to be %s, got %s", node.Key, check.Key)
		}

		if check.Hostname != node.Hostname {
			t.Errorf("Expected Name to be %s, got %s", node.Hostname, check.Hostname)
		}

		if check.Port != node.Port {
			t.Errorf("Expected Port to be %d, got %d", node.Port, check.Port)
		}

		if check.ShardID != node.ShardID {
			t.Errorf("Expected ShardID to be %d, got %d", node.ShardID, check.ShardID)
		}

		if check.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set, got zero value")
		}
	})
}

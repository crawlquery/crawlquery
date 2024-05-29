package service_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/factory"
	"crawlquery/node/dto"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"errors"
	"fmt"
	"testing"
	"time"

	nodeRepo "crawlquery/api/node/repository/mem"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	"crawlquery/api/node/service"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	t.Run("can create a node", func(t *testing.T) {

		accountID := util.UUIDString()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		shardSvc, _ := factory.ShardServiceWithShard(&domain.Shard{
			ID: 1,
		})

		nodeRepo := nodeRepo.NewRepository()

		svc := service.NewService(nodeRepo, accSvc, shardSvc, testutil.NewTestLogger())

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

		if node.Key == "" {
			t.Errorf("Expected Key to be set")
		}

		list, err := nodeRepo.List()

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
		accountID := util.UUIDString()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, nodeRepo := factory.NodeService(accSvc)

		err := nodeRepo.Create(&domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "hostname",
			Port:      8080,
		})

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

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
		accountID := util.UUIDString()
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

		accountID := util.UUIDString()
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

		accountID := util.UUIDString()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		shardSvc, _ := factory.ShardServiceWithShard(&domain.Shard{
			ID: 1,
		})

		nodeRepo := nodeRepo.NewRepository()
		svc := service.NewService(nodeRepo, accSvc, shardSvc, testutil.NewTestLogger())

		nodeRepo.ForceCreateError(errors.New("db locked"))

		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "testnode",
		}

		nodeRepo.Create(node)

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

		accountID := util.UUIDString()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, nodeRepo := factory.NodeService(accSvc)

		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "testnode",
			Port:      8080,
		}

		node2 := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "testnode2",
			Port:      8081,
		}

		nodeRepo.Create(node)
		nodeRepo.Create(node2)

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

		accountID := util.UUIDString()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, repo := factory.NodeService(accSvc)

		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "testnode",
		}

		node2 := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "testnode2",
		}

		node3 := &domain.Node{
			ID:        util.UUIDString(),
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

		accountID := util.UUIDString()
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

func TestGetShardWithLeastNodes(t *testing.T) {
	t.Run("can get shard with least nodes", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		shardRepo := shardRepo.NewRepository()
		shardService := shardService.NewService(shardRepo, testutil.NewTestLogger())
		nodeService := service.NewService(nodeRepo, nil, shardService, testutil.NewTestLogger())

		shard := &domain.Shard{
			ID: 1,
		}

		shard2 := &domain.Shard{
			ID: 2,
		}

		shard3 := &domain.Shard{
			ID: 3,
		}

		shardRepo.Create(shard)
		shardRepo.Create(shard2)
		shardRepo.Create(shard3)

		nodes := []*domain.Node{
			{ID: "1", ShardID: 1},
			{ID: "2", ShardID: 1},
			{ID: "3", ShardID: 2},
			{ID: "4", ShardID: 2},
			{ID: "5", ShardID: 3},
		}

		for _, n := range nodes {
			nodeRepo.Create(n)
		}

		found, err := nodeService.GetShardWithLeastNodes()

		if err != nil {
			t.Fatalf("Error getting shard with least nodes: %v", err)
		}

		if found.ID != 3 {
			t.Errorf("Expected shard ID to be 3, got %d", shard.ID)
		}
	})

	t.Run("can get shard with least nodes when all shards have the same number of nodes", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		shardRepo := shardRepo.NewRepository()
		shardService := shardService.NewService(shardRepo, testutil.NewTestLogger())
		nodeService := service.NewService(nodeRepo, nil, shardService, testutil.NewTestLogger())

		shard := &domain.Shard{
			ID: 1,
		}

		shard2 := &domain.Shard{
			ID: 2,
		}

		shard3 := &domain.Shard{
			ID: 3,
		}

		shardRepo.Create(shard)
		shardRepo.Create(shard2)
		shardRepo.Create(shard3)

		nodes := []*domain.Node{
			{ID: "1", ShardID: 1},
			{ID: "2", ShardID: 1},
			{ID: "3", ShardID: 2},
			{ID: "4", ShardID: 2},
			{ID: "5", ShardID: 3},
			{ID: "6", ShardID: 3},
		}

		for _, n := range nodes {
			nodeRepo.Create(n)
		}

		shard, err := nodeService.GetShardWithLeastNodes()

		if err != nil {
			t.Fatalf("Error getting shard with least nodes: %v", err)
		}

		if shard == nil {
			t.Fatalf("Expected shard to not be nil")
		}
	})

	t.Run("can get shard with least nodes when no nodes exist", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		shardRepo := shardRepo.NewRepository()
		shardService := shardService.NewService(shardRepo, testutil.NewTestLogger())
		nodeService := service.NewService(nodeRepo, nil, shardService, testutil.NewTestLogger())

		shard := &domain.Shard{
			ID: 1,
		}

		shard2 := &domain.Shard{
			ID: 2,
		}

		shard3 := &domain.Shard{
			ID: 3,
		}

		shardRepo.Create(shard)
		shardRepo.Create(shard2)
		shardRepo.Create(shard3)

		shard, err := nodeService.GetShardWithLeastNodes()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if shard == nil {
			t.Errorf("Expected shard to not be nil")
		}
	})
}

func TestAllocateNode(t *testing.T) {
	t.Run("can allocate a node", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		shardRepo := shardRepo.NewRepository()
		shardService := shardService.NewService(shardRepo, testutil.NewTestLogger())
		nodeService := service.NewService(nodeRepo, nil, shardService, testutil.NewTestLogger())

		shard := &domain.Shard{
			ID: 1,
		}

		shardRepo.Create(shard)

		shard2 := &domain.Shard{
			ID: 2,
		}

		shardRepo.Create(shard2)

		nodeRepo.Create(&domain.Node{
			ID:        "1",
			ShardID:   1,
			Hostname:  "testnode",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		nodeRepo.Create(&domain.Node{
			ID:        "2",
			ShardID:   1,
			Hostname:  "testnode2",
			Port:      8080,
			CreatedAt: time.Now(),
		})

		node := &domain.Node{
			ID:        "3",
			Hostname:  "testnode3",
			Port:      8080,
			CreatedAt: time.Now(),
		}

		err := nodeService.AllocateNode(node)

		if err != nil {
			t.Fatalf("Error allocating node: %v", err)
		}

		if node.ShardID != shard2.ID {
			t.Errorf("Expected ShardID to be %d, got %d", shard2.ID, node.ShardID)
		}
	})

	t.Run("can allocate a node when no nodes exist", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		shardRepo := shardRepo.NewRepository()
		shardService := shardService.NewService(shardRepo, testutil.NewTestLogger())
		nodeService := service.NewService(nodeRepo, nil, shardService, testutil.NewTestLogger())

		shard := &domain.Shard{
			ID: 1,
		}

		shardRepo.Create(shard)

		shard2 := &domain.Shard{
			ID: 2,
		}

		shardRepo.Create(shard2)

		node := &domain.Node{
			ID:        "1",
			Hostname:  "testnode",
			Port:      8080,
			CreatedAt: time.Now(),
		}

		err := nodeService.AllocateNode(node)

		if err != nil {
			t.Fatalf("Error allocating node: %v", err)
		}

		if node.ShardID != shard.ID && node.ShardID != shard2.ID {
			t.Errorf("Expected ShardID to be %d or %d, got %d", shard.ID, shard2.ID, node.ShardID)
		}
	})
}

func TestListGroupByShard(t *testing.T) {
	t.Run("can group nodes by shard", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		nodeService := service.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		nodes := []*domain.Node{
			{ID: "1", ShardID: 1},
			{ID: "2", ShardID: 1},
			{ID: "3", ShardID: 2},
			{ID: "4", ShardID: 2},
			{ID: "5", ShardID: 3},
		}

		for _, n := range nodes {
			nodeRepo.Create(n)
		}

		grouped, err := nodeService.ListGroupByShard()

		if err != nil {
			t.Fatalf("Error grouping nodes by shard: %v", err)
		}

		if len(grouped) != 3 {
			t.Fatalf("Expected 3 groups, got %d", len(grouped))
		}

		if len(grouped[1]) != 2 {
			t.Errorf("Expected 2 nodes in shard 1, got %d", len(grouped[1]))
		}

		if len(grouped[2]) != 2 {
			t.Errorf("Expected 2 nodes in shard 2, got %d", len(grouped[2]))
		}

		if len(grouped[3]) != 1 {
			t.Errorf("Expected 1 node in shard 3, got %d", len(grouped[3]))
		}
	})

	t.Run("can group nodes by shard when no nodes exist", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		nodeService := service.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		grouped, err := nodeService.ListGroupByShard()

		if err != nil {
			t.Fatalf("Error grouping nodes by shard: %v", err)
		}

		if len(grouped) != 0 {
			t.Fatalf("Expected 0 groups, got %d", len(grouped))
		}
	})
}

func TestListByAccountID(t *testing.T) {
	t.Run("can list nodes by account ID", func(t *testing.T) {

		accountID := util.UUIDString()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, nodeRepo := factory.NodeService(accSvc)

		node := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "testnode",
			Port:      8080,
		}

		node2 := &domain.Node{
			ID:        util.UUIDString(),
			AccountID: accountID,
			Hostname:  "testnode2",
			Port:      8081,
		}

		nodeRepo.Create(node)
		nodeRepo.Create(node2)

		list, err := svc.ListByAccountID(accountID)

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

	t.Run("can list nodes by account ID when no nodes exist", func(t *testing.T) {

		accountID := util.UUIDString()
		accSvc, _ := factory.AccountServiceWithAccount(&domain.Account{
			ID: accountID,
		})

		svc, _ := factory.NodeService(accSvc)

		list, err := svc.ListByAccountID(accountID)

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(list) != 0 {
			t.Fatalf("Expected 0 nodes, got %d", len(list))
		}
	})
}

func TestListByShardID(t *testing.T) {
	t.Run("can list nodes by shard ID", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		nodeService := service.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		node := &domain.Node{
			ID:      "1",
			ShardID: 1,
		}

		node2 := &domain.Node{
			ID:      "2",
			ShardID: 1,
		}

		node3 := &domain.Node{
			ID:      "3",
			ShardID: 2,
		}

		nodeRepo.Create(node)
		nodeRepo.Create(node2)
		nodeRepo.Create(node3)

		list, err := nodeService.ListByShardID(1)

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(list) != 2 {
			t.Fatalf("Expected 2 nodes, got %d", len(list))
		}

		for _, n := range list {
			if n.ID != "1" && n.ID != "2" {
				t.Errorf("Expected node to be one of 1 or 2, got %s", n.ID)
			}
		}
	})

	t.Run("can list nodes by shard ID when no nodes exist", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		nodeService := service.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		list, err := nodeService.ListByShardID(1)

		if err != nil {
			t.Fatalf("Error listing nodes: %v", err)
		}

		if len(list) != 0 {
			t.Fatalf("Expected 0 nodes, got %d", len(list))
		}
	})
}

func TestSendCrawlJob(t *testing.T) {
	t.Run("can send a crawl job", func(t *testing.T) {
		node := &domain.Node{
			ID:       "1",
			Hostname: "testnode",
			Port:     8080,
		}

		crawlJob := &domain.CrawlJob{
			ID:     "1",
			PageID: "1",
			URL:    "http://google.com",
		}

		defer gock.Off()

		responseJson := fmt.Sprintf(`{"page_id":"%s","url":"%s"}`, crawlJob.PageID, crawlJob.URL)

		crawlResponse := &dto.CrawlResponse{
			Page: &dto.Page{
				ID:   "1",
				URL:  "http://google.com",
				Hash: "hash",
			},
		}

		gock.New("http://testnode:8080").
			Post("/crawl").
			JSON(responseJson).
			Reply(200).
			JSON(crawlResponse)

		nodeService := service.NewService(nil, nil, nil, testutil.NewTestLogger())

		crawledPage, err := nodeService.SendCrawlJob(node, crawlJob)

		if err != nil {
			t.Fatalf("Error sending crawl job: %v", err)
		}

		if crawledPage.Hash != crawlResponse.Page.Hash {
			t.Errorf("Expected page hash to be %s, got %s", crawlResponse.Page.Hash, crawledPage.Hash)
		}
	})

	t.Run("handles error sending crawl job", func(t *testing.T) {
		node := &domain.Node{
			ID:       "1",
			Hostname: "testnode",
			Port:     8080,
		}

		crawlJob := &domain.CrawlJob{
			ID:  "1",
			URL: "http://google.com",
		}

		defer gock.Off()

		gock.New("http://testnode:8080").
			Post("/crawl").
			JSON(`{"url":"http://google.com"`).
			Reply(500)

		nodeService := service.NewService(nil, nil, nil, testutil.NewTestLogger())

		_, err := nodeService.SendCrawlJob(node, crawlJob)

		if err == nil {
			t.Fatalf("Expected error sending crawl job")
		}
	})
}

func TestSendIndexJob(t *testing.T) {
	t.Run("can send an index job", func(t *testing.T) {
		node := &domain.Node{
			ID:       "1",
			Hostname: "testnode",
			Port:     8080,
		}

		indexJob := &domain.IndexJob{
			ID:     "1",
			PageID: "1",
		}

		defer gock.Off()

		indexResponse := &dto.IndexResponse{
			Success: true,
		}

		gock.New("http://testnode:8080").
			Post(fmt.Sprintf("/pages/%s/index", indexJob.PageID)).
			Reply(200).
			JSON(indexResponse)

		nodeService := service.NewService(nil, nil, nil, testutil.NewTestLogger())

		err := nodeService.SendIndexJob(node, indexJob)

		if err != nil {
			t.Fatalf("Error sending index job: %v", err)
		}
	})

	t.Run("handles error sending index job", func(t *testing.T) {
		node := &domain.Node{
			ID:       "1",
			Hostname: "testnode",
			Port:     8080,
		}

		indexJob := &domain.IndexJob{
			ID:     "1",
			PageID: "1",
		}

		defer gock.Off()

		gock.New("http://testnode:8080").
			Post("/pages/1/index").
			Reply(500)

		nodeService := service.NewService(nil, nil, nil, testutil.NewTestLogger())

		err := nodeService.SendIndexJob(node, indexJob)

		if err == nil {
			t.Fatalf("Expected error sending index job")
		}
	})
}

func TestRandomize(t *testing.T) {
	t.Run("can randomize a list of nodes", func(t *testing.T) {
		nodeRepo := nodeRepo.NewRepository()
		nodeService := service.NewService(nodeRepo, nil, nil, testutil.NewTestLogger())

		nodes := []*domain.Node{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
			{ID: "4"},
			{ID: "5"},
		}

		for _, n := range nodes {
			nodeRepo.Create(n)
		}

		randomized := nodeService.Randomize(nodes)

		if len(randomized) != 5 {
			t.Fatalf("Expected 5 nodes, got %d", len(randomized))
		}

		if randomized[0].ID == "1" && randomized[1].ID == "2" && randomized[2].ID == "3" && randomized[3].ID == "4" && randomized[4].ID == "5" {
			t.Fatalf("Expected nodes to be randomized")
		}
	})
}

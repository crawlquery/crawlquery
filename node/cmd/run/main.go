package main

import (
	crawlHandler "crawlquery/node/crawl/handler"
	crawlService "crawlquery/node/crawl/service"
	"crawlquery/node/domain"
	htmlRepo "crawlquery/node/html/repository/disk"
	htmlService "crawlquery/node/html/service"
	"crawlquery/pkg/client/api"
	"fmt"
	"time"

	pageRepo "crawlquery/node/page/repository/bolt"
	pageService "crawlquery/node/page/service"

	peerService "crawlquery/node/peer/service"

	indexHandler "crawlquery/node/index/handler"
	indexService "crawlquery/node/index/service"

	dumpHandler "crawlquery/node/dump/handler"
	dumpService "crawlquery/node/dump/service"

	"crawlquery/node/router"

	"flag"

	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	var key string
	var apiURL string

	flag.StringVar(&key, "key", "secret", "")
	flag.StringVar(&apiURL, "api", "http://localhost:8080", "API BaseURL")

	var htmlStoragePath string
	var pageDBPath string

	flag.StringVar(&htmlStoragePath, "html", "/tmp/htmlstorage", "path to the html storage")
	flag.StringVar(&pageDBPath, "pdb", "/tmp/pagedb.bolt", "path to the pagedb")

	flag.Parse()

	// Create repositories
	htmlRepo, err := htmlRepo.NewRepository(htmlStoragePath)
	if err != nil {
		sugar.Fatalf("Error creating html repository: %v", err)
	}

	pageRepo, err := pageRepo.NewRepository(pageDBPath)
	if err != nil {
		sugar.Fatalf("Error creating page repository: %v", err)
	}

	api := api.NewClient(apiURL, sugar)

	// Authenticate with the API
	node, err := api.AuthenticateNode(key)

	if err != nil {
		sugar.Fatalf("Error authenticating node: %v", err)
	}

	fmt.Printf("Node ID: %s\n", node.ID)
	fmt.Printf("Node Hostname: %s\n", node.Hostname)
	fmt.Printf("Node Port: %d\n", node.Port)
	fmt.Printf("Node Shard ID: %d\n", node.ShardID)

	// Create services
	htmlService := htmlService.NewService(htmlRepo)
	pageService := pageService.NewService(pageRepo)
	peerService := peerService.NewService(api, pageService, &domain.Peer{
		ID:       node.ID,
		Hostname: node.Hostname,
		Port:     node.Port,
		ShardID:  node.ShardID,
	}, sugar)
	indexService := indexService.NewService(pageService, htmlService, peerService, sugar)
	crawlService := crawlService.NewService(htmlService, pageService, indexService, api, sugar)
	dumpService := dumpService.NewService(pageService)

	// Create handlers
	indexHandler := indexHandler.NewHandler(indexService, sugar)
	crawlHandler := crawlHandler.NewHandler(crawlService, sugar)
	dumpHandler := dumpHandler.NewHandler(dumpService)

	peerService.SyncPeerList()
	go peerService.SyncPeerListEvery(30 * time.Second)

	r := router.NewRouter(indexHandler, crawlHandler, dumpHandler)

	r.Run(fmt.Sprintf(":%d", node.Port))
}

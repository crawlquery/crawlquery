package main

import (
	crawlHandler "crawlquery/node/crawl/handler"
	crawlService "crawlquery/node/crawl/service"
	"crawlquery/node/domain"
	htmlBackupService "crawlquery/node/html/backup/service"
	htmlRepo "crawlquery/node/html/repository/disk"
	htmlService "crawlquery/node/html/service"
	htmlClient "crawlquery/pkg/client/html"

	"crawlquery/pkg/client/api"
	"fmt"
	"time"

	pageRepo "crawlquery/node/page/repository/bolt"
	pageService "crawlquery/node/page/service"

	keywordRepo "crawlquery/node/keyword/repository/bolt"
	keywordService "crawlquery/node/keyword/service"

	peerService "crawlquery/node/peer/service"

	indexHandler "crawlquery/node/index/handler"
	indexService "crawlquery/node/index/service"

	dumpHandler "crawlquery/node/dump/handler"
	dumpService "crawlquery/node/dump/service"

	statHandler "crawlquery/node/stat/handler"
	statService "crawlquery/node/stat/service"

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
	var htmlBackupURL string
	var pageDBPath string
	var keywordDBPath string

	flag.StringVar(&htmlStoragePath, "html", "/tmp/htmlstorage", "path to the html storage")
	flag.StringVar(&pageDBPath, "pdb", "/tmp/pagedb.bolt", "path to the pagedb")
	flag.StringVar(&keywordDBPath, "kdb", "/tmp/keyworddb.bolt", "path to the keyworddb")
	flag.StringVar(&htmlBackupURL, "htmlbackup", "http://crawlquery-html1.dxs.network", "URL to the html backup service")

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

	keywordRepo, err := keywordRepo.NewRepository(keywordDBPath)
	if err != nil {
		sugar.Fatalf("Error creating keyword repository: %v", err)
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

	// clients
	htmlClient := htmlClient.NewClient(htmlBackupURL)

	// services
	htmlBackupService := htmlBackupService.NewService(htmlClient)
	htmlService := htmlService.NewService(htmlRepo, htmlBackupService)
	peerService := peerService.NewService(api, &domain.Peer{
		ID:       node.ID,
		Hostname: node.Hostname,
		Port:     node.Port,
		ShardID:  node.ShardID,
	}, sugar)
	pageService := pageService.NewService(pageRepo, peerService)
	keywordService := keywordService.NewService(keywordRepo)
	indexService := indexService.NewService(pageService, htmlService, peerService, keywordService, sugar)
	crawlService := crawlService.NewService(htmlService, pageService, indexService, api, sugar)
	dumpService := dumpService.NewService(pageService)
	statService := statService.NewService(pageService, dumpService)

	// handlers
	indexHandler := indexHandler.NewHandler(indexService, sugar)
	crawlHandler := crawlHandler.NewHandler(crawlService, sugar)
	dumpHandler := dumpHandler.NewHandler(dumpService)
	statHandler := statHandler.NewHandler(statService)

	peerService.SyncPeerList()
	go peerService.SyncPeerListEvery(30 * time.Second)

	r := router.NewRouter(indexHandler, crawlHandler, dumpHandler, statHandler)

	r.Run(fmt.Sprintf(":%d", node.Port))
}

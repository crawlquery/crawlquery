package main

import (
	crawlHandler "crawlquery/node/crawl/handler"
	crawlService "crawlquery/node/crawl/service"
	htmlRepo "crawlquery/node/html/repository/disk"
	htmlService "crawlquery/node/html/service"

	pageRepo "crawlquery/node/page/repository/bolt"
	pageService "crawlquery/node/page/service"

	keywordRepo "crawlquery/node/keyword/repository/bolt"
	keywordService "crawlquery/node/keyword/service"

	indexHandler "crawlquery/node/index/handler"
	indexService "crawlquery/node/index/service"
	"crawlquery/node/router"

	"flag"

	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	var portFlag string
	var htmlStoragePath string
	var pageDBPath string
	var keywordDBPath string

	flag.StringVar(&portFlag, "port", "9090", "port to run the server on")
	flag.StringVar(&htmlStoragePath, "htmlstorage", "/tmp/htmlstorage", "path to the html storage")
	flag.StringVar(&pageDBPath, "pagedb", "/tmp/pagedb.bolt", "path to the pagedb")
	flag.StringVar(&keywordDBPath, "keyworddb", "/tmp/keyworddb.bolt", "path to the keyworddb")

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

	// Create services
	htmlService := htmlService.NewService(htmlRepo)
	pageService := pageService.NewService(pageRepo)
	keywordService := keywordService.NewService(keywordRepo)
	indexService := indexService.NewService(pageService, htmlService, keywordService, sugar)
	crawlService := crawlService.NewService(htmlService, pageService, indexService, sugar)

	// Create handlers
	indexHandler := indexHandler.NewHandler(indexService, sugar)
	crawlHandler := crawlHandler.NewHandler(crawlService, sugar)

	r := router.NewRouter(indexHandler, crawlHandler)

	r.Run(":" + portFlag)
}

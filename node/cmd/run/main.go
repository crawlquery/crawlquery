package main

import (
	crawlHandler "crawlquery/node/crawl/handler"
	crawlService "crawlquery/node/crawl/service"
	htmlRepo "crawlquery/node/html/repository/disk"
	"crawlquery/node/index"
	indexHandler "crawlquery/node/index/handler"
	"crawlquery/node/router"

	forwardRepo "crawlquery/node/index/forward/repository/bolt"
	invertedRepo "crawlquery/node/index/inverted/repository/bolt"
	"flag"

	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	var portFlag string
	var forwardIdxFlag string
	var invertedIdxFlag string

	flag.StringVar(&portFlag, "port", "9090", "port to run the server on")
	flag.StringVar(&forwardIdxFlag, "forwardIdx", "/tmp/forwardidx.bolt", "path to the forward index")
	flag.StringVar(&invertedIdxFlag, "invertedIdx", "/tmp/invertedidx.bolt", "path to the inverted index")

	flag.Parse()

	forwardRepo, err := forwardRepo.NewRepository(forwardIdxFlag)

	if err != nil {
		panic(err)
	}

	invertedRepo, err := invertedRepo.NewRepository(invertedIdxFlag)

	if err != nil {
		panic(err)
	}

	idx := index.NewIndex(
		forwardRepo,
		invertedRepo,
		sugar,
	)

	htmlRepo, err := htmlRepo.NewRepository("/tmp/crawlquery-html")

	if err != nil {
		panic(err)
	}

	indexHandler := indexHandler.NewHandler(idx)
	crawlHandler := crawlHandler.NewHandler(
		crawlService.NewService(
			htmlRepo,
			sugar,
		),
	)

	r := router.NewRouter(indexHandler, crawlHandler)

	r.Run(":" + portFlag)
}

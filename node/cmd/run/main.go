package main

import (
	crawlHandler "crawlquery/node/crawl/handler"
	crawlService "crawlquery/node/crawl/service"
	"crawlquery/node/index"
	indexHandler "crawlquery/node/index/handler"
	"crawlquery/node/router"
	"crawlquery/pkg/factory"
	"flag"
)

func main() {

	var portFlag string

	flag.StringVar(&portFlag, "port", "9090", "port to run the server on")

	flag.Parse()

	idx := index.NewIndex()
	for _, page := range factory.TenPages() {
		idx.AddPage(page)
	}

	indexHandler := indexHandler.NewHandler(idx)
	crawlHandler := crawlHandler.NewHandler(
		crawlService.NewService(
			idx,
		),
	)

	r := router.NewRouter(indexHandler, crawlHandler)

	r.Run(":" + portFlag)
}

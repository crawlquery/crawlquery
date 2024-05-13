package main

import (
	crawlHandler "crawlquery/node/crawl/handler"
	crawlService "crawlquery/node/crawl/service"
	"crawlquery/node/index"
	indexHandler "crawlquery/node/index/handler"
	"crawlquery/node/router"
	"crawlquery/pkg/factory"
)

func main() {

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

	r.Run(":9090")
}

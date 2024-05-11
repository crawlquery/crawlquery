package main

import (
	"crawlquery/api/handler"
	"crawlquery/api/router"
	"crawlquery/api/service"
	"crawlquery/pkg/repository/node/disk"
)

func main() {

	searchService := service.NewSearchService(
		service.NewNodeService(
			disk.NewDiskRepository(
				"/tmp/nodes.gob",
			),
		),
	)
	searchHandler := handler.NewSearchHandler(searchService)
	r := router.NewRouter(searchHandler)

	r.Run(":8080")
}

package main

import (
	"crawlquery/api/handler"
	"crawlquery/api/router"
	"crawlquery/api/service"
	crawlJobDiskRepo "crawlquery/pkg/repository/job/disk"
	nodeDiskRepo "crawlquery/pkg/repository/node/disk"
)

func main() {

	cqr := crawlJobDiskRepo.NewDiskRepository(
		"/tmp/crawl_job.gob",
	)

	ns := service.NewNodeService(
		nodeDiskRepo.NewDiskRepository(
			"/tmp/nodes.gob",
		),
	)
	searchService := service.NewSearchService(
		ns,
	)
	searchHandler := handler.NewSearchHandler(searchService)

	crawlHandler := handler.NewCrawlHandler(
		service.NewCrawlService(
			cqr,
			ns,
		),
	)

	r := router.NewRouter(searchHandler, crawlHandler)

	r.Run(":8080")
}

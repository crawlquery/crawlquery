package main

import (
	"crawlquery/api/handler"
	"crawlquery/api/router"
	"crawlquery/api/service"
	crawlJobDiskRepo "crawlquery/pkg/repository/job/disk"
	nodeDiskRepo "crawlquery/pkg/repository/node/disk"
	"fmt"
	"time"
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
	crawlSvc := service.NewCrawlService(
		cqr,
		ns,
	)
	crawlHandler := handler.NewCrawlHandler(
		crawlSvc,
	)

	nodeHandler := handler.NewNodeHandler(ns)

	r := router.NewRouter(
		searchHandler,
		crawlHandler,
		nodeHandler,
	)

	go func() {
		for {
			fmt.Println("Crawling...")
			err := crawlSvc.Crawl()

			if err != nil {
				fmt.Println("Error crawling: ", err)
			}

			time.Sleep(time.Second)
		}
	}()

	r.Run(":8080")
}

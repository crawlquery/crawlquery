package main

import (
	"crawlquery/node/handler"
	"crawlquery/node/router"
	"crawlquery/node/service"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/index"
	"crawlquery/pkg/repository/index/mem"
)

func main() {

	idx := index.NewIndex()

	for _, page := range factory.TenPages() {
		idx.AddPage(page)
	}

	memRepo := mem.NewMemoryRepository()
	memRepo.Save(idx)

	svc := service.NewIndexService(memRepo)
	handler := handler.NewSearchHandler(svc)

	r := router.NewRouter(handler)

	r.Run(":9090")
}

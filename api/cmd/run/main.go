package main

import (
	"crawlquery/api/migration"
	"crawlquery/api/router"
	"database/sql"
	"fmt"
	"os"
	"time"

	authHandler "crawlquery/api/auth/handler"
	authService "crawlquery/api/auth/service"

	accountHandler "crawlquery/api/account/handler"
	accountMysqlRepo "crawlquery/api/account/repository/mysql"
	accountService "crawlquery/api/account/service"

	crawlJobMysqlRepo "crawlquery/api/crawl/job/repository/mysql"
	crawlLogMysqlRepo "crawlquery/api/crawl/log/repository/mysql"
	crawlService "crawlquery/api/crawl/service"

	nodeHandler "crawlquery/api/node/handler"
	nodeMysqlRepo "crawlquery/api/node/repository/mysql"
	nodeService "crawlquery/api/node/service"

	shardMysqlRepo "crawlquery/api/shard/repository/mysql"
	shardService "crawlquery/api/shard/service"

	searchHandler "crawlquery/api/search/handler"
	searchService "crawlquery/api/search/service"

	linkHandler "crawlquery/api/link/handler"
	linkMySQLRepo "crawlquery/api/link/repository/mysql"
	linkService "crawlquery/api/link/service"

	// pageRankMysqlRepo "crawlquery/api/pagerank/repository/mysql"
	pageRankMemRepo "crawlquery/api/pagerank/repository/mem"
	pageRankService "crawlquery/api/pagerank/service"

	pageMysqlRepo "crawlquery/api/page/repository/mysql"
	pageService "crawlquery/api/page/service"

	indexJobMySQLRepo "crawlquery/api/index/job/repository/mysql"
	indexJobService "crawlquery/api/index/job/service"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := sql.Open("mysql", "root:cqdb@tcp(localhost:3306)/cqdb_v2?parseTime=true")
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}
	defer db.Close()

	if os.Getenv("ENV") == "development" {
		fmt.Println("Running migrations...")
		err := migration.Up(db)

		if err != nil {
			fmt.Println("Error migrating database: ", err)
			return
		}
	}

	accountRepo := accountMysqlRepo.NewRepository(db)
	accountService := accountService.NewService(accountRepo, sugar)
	accountHandler := accountHandler.NewHandler(accountService)

	authService := authService.NewService(accountService, sugar)
	authHandler := authHandler.NewHandler(authService)

	shardRepo := shardMysqlRepo.NewRepository(db)
	shardService := shardService.NewService(
		shardService.WithRepo(shardRepo),
		shardService.WithLogger(sugar),
	)

	nodeRepo := nodeMysqlRepo.NewRepository(db)
	nodeService := nodeService.NewService(nodeRepo, accountService, shardService, sugar)
	nodeHandler := nodeHandler.NewHandler(nodeService)

	pageRepo := pageMysqlRepo.NewRepository(db)
	pageService := pageService.NewService(
		pageService.WithPageRepo(pageRepo),
		pageService.WithShardService(shardService),
		pageService.WithLogger(sugar),
	)
	pageHandler := pageHandler.NewHandler(pageService)

	indexJobRepo := indexJobMySQLRepo.NewRepository(db)
	indexJobService := indexJobService.NewService(indexJobRepo, pageService, nodeService, sugar)

	crawlJobRepo := crawlJobMysqlRepo.NewRepository(db)
	crawlJobService := crawlService.NewService(
		crawlService.WithCrawlJobRepo(crawlJobRepo),
		crawlService.WithCrawlLogRepo(crawlLogMysqlRepo.NewRepository(db)),
		crawlService.WithLogger(sugar),
	)

	linkRepo := linkMySQLRepo.NewRepository(db)
	linkService := linkService.NewService(linkRepo, crawlJobService, sugar)
	linkHandler := linkHandler.NewHandler(linkService, sugar)

	// pageRankRepo := pageRankMysqlRepo.NewRepository(db)
	pageRankRepo := pageRankMemRepo.NewRepository()
	pageRankService := pageRankService.NewService(linkService, pageRankRepo, sugar)

	searchService := searchService.NewService(nodeService, pageRankService, sugar)
	searchHandler := searchHandler.NewHandler(searchService)

	go crawlJobService.ProcessCrawlJobs()

	go pageRankService.UpdatePageRanksEvery(time.Minute)

	// start 4 workers to process index jobs
	for i := 0; i < 4; i++ {
		go indexJobService.ProcessIndexJobs()
	}

	r := router.NewRouter(
		accountService,
		authHandler,
		accountHandler,
		pageHandler,
		nodeHandler,
		searchHandler,
		linkHandler,
	)

	r.Run(":8080")
}

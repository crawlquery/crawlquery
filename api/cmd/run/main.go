package main

import (
	"context"
	"crawlquery/api/migration"
	"crawlquery/api/router"
	"database/sql"
	"fmt"
	"os"
	"time"

	eventService "crawlquery/api/event/service"

	authHandler "crawlquery/api/auth/handler"
	authService "crawlquery/api/auth/service"

	accountHandler "crawlquery/api/account/handler"
	accountMysqlRepo "crawlquery/api/account/repository/mysql"
	accountService "crawlquery/api/account/service"

	crawlJobMysqlRepo "crawlquery/api/crawl/job/repository/mysql"
	crawlLogMysqlRepo "crawlquery/api/crawl/log/repository/mysql"
	crawlService "crawlquery/api/crawl/service"
	crawlThrottleService "crawlquery/api/crawl/throttle/service"

	nodeHandler "crawlquery/api/node/handler"
	nodeMysqlRepo "crawlquery/api/node/repository/mysql"
	nodeService "crawlquery/api/node/service"

	shardMysqlRepo "crawlquery/api/shard/repository/mysql"
	shardService "crawlquery/api/shard/service"

	searchHandler "crawlquery/api/search/handler"
	searchService "crawlquery/api/search/service"

	linkMySQLRepo "crawlquery/api/link/repository/mysql"
	linkService "crawlquery/api/link/service"

	// pageRankMysqlRepo "crawlquery/api/pagerank/repository/mysql"
	pageRankMemRepo "crawlquery/api/pagerank/repository/mem"
	pageRankService "crawlquery/api/pagerank/service"

	pageHandler "crawlquery/api/page/handler"
	pageMysqlRepo "crawlquery/api/page/repository/mysql"
	pageService "crawlquery/api/page/service"

	indexJobMySQLRepo "crawlquery/api/index/job/repository/mysql"
	indexLogMysqlRepo "crawlquery/api/index/log/repository/mysql"
	indexService "crawlquery/api/index/service"

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

	eventService := eventService.NewService()

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
	nodeService := nodeService.NewService(
		nodeService.WithAccountService(accountService),
		nodeService.WithNodeRepo(nodeRepo),
		nodeService.WithShardService(shardService),
		nodeService.WithLogger(sugar),
	)
	nodeHandler := nodeHandler.NewHandler(nodeService)

	pageRepo := pageMysqlRepo.NewRepository(db)
	pageService := pageService.NewService(
		pageService.WithPageRepo(pageRepo),
		pageService.WithShardService(shardService),
		pageService.WithLogger(sugar),
		pageService.WithEventService(eventService),
	)
	pageHandler := pageHandler.NewHandler(pageService)

	indexJobRepo := indexJobMySQLRepo.NewRepository(db)
	indexLogRepo := indexLogMysqlRepo.NewRepository(db)
	indexService := indexService.NewService(
		indexService.WithEventService(eventService),
		indexService.WithEventListeners(),
		indexService.WithIndexJobRepo(indexJobRepo),
		indexService.WithIndexLogRepo(indexLogRepo),
		indexService.WithNodeService(nodeService),
		indexService.WithLogger(sugar),
		indexService.WithWorkers(4),
		indexService.WithMaxQueueSize(100),
	)

	crawlThrottleService := crawlThrottleService.NewService(
		crawlThrottleService.WithRateLimit(time.Second * 20),
	)
	crawlJobRepo := crawlJobMysqlRepo.NewRepository(db)
	crawlJobService := crawlService.NewService(
		crawlService.WithEventService(eventService),
		crawlService.WithEventListeners(),
		crawlService.WithCrawlThrottleService(crawlThrottleService),
		crawlService.WithCrawlJobRepo(crawlJobRepo),
		crawlService.WithCrawlLogRepo(crawlLogMysqlRepo.NewRepository(db)),
		crawlService.WithLogger(sugar),
		crawlService.WithWorkers(10),
		crawlService.WithMaxQueueSize(10000),
	)

	linkRepo := linkMySQLRepo.NewRepository(db)
	linkService := linkService.NewService(
		linkService.WithLinkRepo(linkRepo),
		linkService.WithLogger(sugar),
		linkService.WithEventService(eventService),
		linkService.WithEventListeners(),
	)

	// pageRankRepo := pageRankMysqlRepo.NewRepository(db)
	pageRankRepo := pageRankMemRepo.NewRepository()
	pageRankService := pageRankService.NewService(linkService, pageRankRepo, sugar)

	searchService := searchService.NewService(nodeService, pageRankService, sugar)
	searchHandler := searchHandler.NewHandler(searchService)

	go crawlJobService.RunCrawlProcess(context.Background())

	go pageRankService.UpdatePageRanksEvery(time.Minute)

	go indexService.RunIndexProcess(context.Background())

	r := router.NewRouter(
		accountService,
		authHandler,
		accountHandler,
		pageHandler,
		nodeHandler,
		searchHandler,
	)

	r.Run(":8080")
}

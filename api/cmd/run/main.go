package main

import (
	"crawlquery/api/migration"
	"crawlquery/api/router"
	"database/sql"
	"fmt"
	"os"

	accountHandler "crawlquery/api/account/handler"
	accountMysqlRepo "crawlquery/api/account/repository/mysql"
	accountService "crawlquery/api/account/service"

	crawlHandler "crawlquery/api/crawl/job/handler"
	crawlJobMysqlRepo "crawlquery/api/crawl/job/repository/mysql"
	crawlJobService "crawlquery/api/crawl/job/service"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := sql.Open("mysql", "root:cqdb@tcp(localhost:3306)/cqdb")
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}
	defer db.Close()

	if os.Getenv("ENV") == "development" {
		err := migration.Up(db)

		if err != nil {
			fmt.Println("Error migrating database: ", err)
			return
		}
	}

	accountRepo := accountMysqlRepo.NewRepository(db)
	accountService := accountService.NewService(accountRepo, sugar)
	accountHandler := accountHandler.NewHandler(accountService)

	crawlJobRepo := crawlJobMysqlRepo.NewRepository(db)
	crawlJobService := crawlJobService.NewService(crawlJobRepo, sugar)
	crawlJobHandler := crawlHandler.NewHandler(crawlJobService)

	r := router.NewRouter(accountHandler, crawlJobHandler)

	r.Run(":8080")
}

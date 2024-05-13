package router

import (
	"crawlquery/api/domain"
	"crawlquery/api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	as domain.AccountService,
	accountHandler domain.AccountHandler,
	crawlJobHandler domain.CrawlJobHandler,
	nodeHandler domain.NodeHandler,
) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("/accounts", accountHandler.Create)
	router.POST("/crawl-jobs", crawlJobHandler.Create)

	router.POST("/nodes", middleware.AuthMiddleware(as, nodeHandler.Create))
	return router
}

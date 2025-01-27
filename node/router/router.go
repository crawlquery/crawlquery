package router

import (
	"crawlquery/node/domain"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewRouter(
	indexHandler domain.IndexHandler,
	searchHandler domain.SearchHandler,
	queryHandler domain.QueryHandler,
	crawlHandler domain.CrawlHandler,
	dumpHandler domain.DumpHandler,
	statHandler domain.StatHandler,
	repairHandler domain.RepairHandler,
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
	router.GET("/search", searchHandler.Search)
	router.POST("/crawl", crawlHandler.Crawl)

	router.POST("/event", indexHandler.Event)
	router.GET("/hash/index", indexHandler.Hash)

	router.GET("/pages/:pageID/index", indexHandler.GetIndex)
	router.POST("/index", indexHandler.Index)

	router.POST("/query", queryHandler.Query)

	router.GET("/dump/page", dumpHandler.Page)

	router.GET("/stats", statHandler.Info)

	router.POST("/repair/get-index-metas", repairHandler.GetIndexMetas)
	router.POST("/repair/get-page-dumps", repairHandler.GetPageDumps)
	router.GET("/repair/get-all-index-metas", repairHandler.GetAllIndexMetas)

	return router
}

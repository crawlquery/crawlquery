package handler

import (
	"crawlquery/api/service"
	"net/url"

	"github.com/gin-gonic/gin"
)

type CrawlHandler struct {
	crawlService *service.CrawlService
}

func NewCrawlHandler(cs *service.CrawlService) *CrawlHandler {
	return &CrawlHandler{
		crawlService: cs,
	}
}

func (ch *CrawlHandler) Crawl(c *gin.Context) {
	if c.Query("url") == "" {
		c.JSON(400, gin.H{
			"error": "url is required",
		})
		return
	}

	// check url is valid
	url, err := url.ParseRequestURI(c.Query("url"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "url is invalid",
		})
		return
	}

	err = ch.crawlService.Queue(url.String())

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}

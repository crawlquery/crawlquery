package handler

import (
	"crawlquery/node/domain"
	"net/url"

	"github.com/gin-gonic/gin"
)

type CrawlHandler struct {
	crawlService domain.CrawlService
}

func NewHandler(cs domain.CrawlService) *CrawlHandler {
	return &CrawlHandler{
		crawlService: cs,
	}
}

func (ch *CrawlHandler) Crawl(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.URL == "" {
		c.JSON(400, gin.H{
			"error": "url is required",
		})
		return
	}

	// check url is valid
	_, err := url.ParseRequestURI(req.URL)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "url is invalid",
		})
		return
	}

	err = ch.crawlService.Crawl(req.URL)

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

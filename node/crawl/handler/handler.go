package handler

import (
	"crawlquery/node/domain"
	"crawlquery/pkg/dto"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CrawlHandler struct {
	crawlService domain.CrawlService
	logger       *zap.SugaredLogger
}

func NewHandler(
	cs domain.CrawlService,
	logger *zap.SugaredLogger,
) *CrawlHandler {
	return &CrawlHandler{
		crawlService: cs,
		logger:       logger,
	}
}

func (ch *CrawlHandler) Crawl(c *gin.Context) {

	var req dto.CrawlRequest

	if err := c.BindJSON(&req); err != nil {
		ch.logger.Errorw("Error binding request", "error", err)
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
		ch.logger.Errorw("Error parsing url", "error", err)
		c.JSON(400, gin.H{
			"error": "url is invalid",
		})
		return
	}

	err = ch.crawlService.Crawl(req.PageID, req.URL)

	if err != nil {
		ch.logger.Errorw("Error crawling page", "error", err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ch.logger.Infow("Page crawled", "pageID", req.PageID, "url", req.URL)
	c.JSON(200, gin.H{
		"message": "success",
	})
}

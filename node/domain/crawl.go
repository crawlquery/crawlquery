package domain

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrCrawlFailedToStoreHtml = errors.New("failed to store html")
var ErrCrawlFailedToFetchHtml = errors.New("failed to fetch html")

type CrawlHandler interface {
	Crawl(c *gin.Context)
}

type CrawlService interface {
	Crawl(pageID, url string) (string, []string, error)
}

package service

import (
	"crawlquery/node/domain"
	"crawlquery/node/html/repository/disk"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

type CrawlService struct {
	htmlRepository *disk.Repository
	logger         *zap.SugaredLogger
}

func NewService(htmlRepo *disk.Repository, logger *zap.SugaredLogger) *CrawlService {
	return &CrawlService{
		htmlRepository: htmlRepo,
		logger:         logger,
	}
}

func (cs *CrawlService) Crawl(pageID, url string) error {

	// Instantiate default collector
	c := colly.NewCollector()

	var failedErr error

	c.OnResponse(func(r *colly.Response) {
		err := cs.htmlRepository.Save(pageID, r.Body)
		if err != nil {
			cs.logger.Errorw("Error saving page", "error", err, "pageID", pageID)
			failedErr = domain.ErrCrawlFailedToStoreHtml
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		cs.logger.Errorw("Error crawling page", "error", e, "pageID", pageID)
		failedErr = domain.ErrCrawlFailedToFetchHtml
	})

	c.Visit(url)

	return failedErr
}

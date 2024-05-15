package service

import (
	"crawlquery/node/domain"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

type CrawlService struct {
	htmlService domain.HTMLService
	pageService domain.PageService
	logger      *zap.SugaredLogger
}

func NewService(
	htmlService domain.HTMLService,
	pageService domain.PageService,
	logger *zap.SugaredLogger,
) *CrawlService {
	return &CrawlService{
		htmlService: htmlService,
		pageService: pageService,
		logger:      logger,
	}
}

func (cs *CrawlService) Crawl(pageID, url string) error {

	// Instantiate default collector
	c := colly.NewCollector()

	var failedErr error

	c.OnResponse(func(r *colly.Response) {
		err := cs.htmlService.Save(pageID, r.Body)
		if err != nil {
			cs.logger.Errorw("Error saving page", "error", err, "pageID", pageID)
			failedErr = domain.ErrCrawlFailedToStoreHtml
			return
		}

		if found, err := cs.pageService.Get(pageID); err == nil && found.ID == pageID {
			cs.logger.Info("Page already exists", "pageID", pageID)
			return
		}

		page, err := cs.pageService.Create(pageID, url)

		if err != nil {
			cs.logger.Errorw("Error creating page", "error", err, "pageID", pageID)
			failedErr = err
		}

		cs.logger.Infow("Page created", "pageID", page.ID, "url", page.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		cs.logger.Errorw("Error crawling page", "error", e, "pageID", pageID)
		failedErr = domain.ErrCrawlFailedToFetchHtml
	})

	c.Visit(url)

	return failedErr
}

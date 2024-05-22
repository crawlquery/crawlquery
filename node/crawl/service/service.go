package service

import (
	"crawlquery/node/domain"
	"crawlquery/pkg/client/api"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

type CrawlService struct {
	htmlService  domain.HTMLService
	pageService  domain.PageService
	indexService domain.IndexService
	api          *api.Client
	logger       *zap.SugaredLogger
}

func NewService(
	htmlService domain.HTMLService,
	pageService domain.PageService,
	indexService domain.IndexService,
	api *api.Client,
	logger *zap.SugaredLogger,
) *CrawlService {
	return &CrawlService{
		htmlService:  htmlService,
		pageService:  pageService,
		indexService: indexService,
		api:          api,
		logger:       logger,
	}
}

func (cs *CrawlService) Crawl(pageID, url string) error {

	// Instantiate default collector
	c := colly.NewCollector()

	var failedErr error

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			cs.logger.Errorw("Error fetching page", "status", r.StatusCode, "pageID", pageID)

			failedErr = domain.ErrCrawlFailedToFetchHtml
			return
		}

		err := cs.htmlService.Save(pageID, r.Body)
		if err != nil {
			cs.logger.Errorw("Error saving page", "error", err, "pageID", pageID)
			failedErr = domain.ErrCrawlFailedToStoreHtml
			return
		}
		page, err := cs.pageService.Get(pageID)
		if err != nil && page != nil && page.ID != pageID {
			cs.logger.Info("Existing page found", "pageID", pageID)
		} else {
			page, err = cs.pageService.Create(pageID, url)
			if err != nil {
				cs.logger.Errorw("Error creating page", "error", err, "pageID", pageID)
				failedErr = err
			}
		}

		cs.logger.Infow("Page created", "pageID", page.ID, "url", page.URL)

		if err := cs.indexService.Index(pageID); err != nil {
			cs.logger.Errorw("Error indexing page", "error", err, "pageID", pageID)
			failedErr = domain.ErrCrawlFailedToIndexPage
		}

		cs.logger.Infow("Page indexed", "pageID", page.ID)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if _, err := cs.api.CreateCrawlJob(link); err != nil {
			cs.logger.Errorw("Error crawling link", "error", err, "link", link)
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		cs.logger.Errorw("Error crawling page", "error", e, "pageID", pageID)
		failedErr = e
	})

	c.Visit(url)

	return failedErr
}

package service

import (
	"crawlquery/node/domain"
	"crawlquery/pkg/client/api"
	"crawlquery/pkg/util"
	"strings"

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

func (cs *CrawlService) Crawl(pageID, url string) (*domain.Page, error) {

	// Instantiate default collector
	c := colly.NewCollector()

	var failedErr error
	var pageHash string
	var pageCrawled *domain.Page

	c.OnResponse(func(r *colly.Response) {

		if !strings.Contains(r.Headers.Get("Content-Type"), "text/html") {
			cs.logger.Errorw("Error fetching page", "error", "Content-Type is not text/html", "pageID", pageID)
			failedErr = domain.ErrCrawlFailedToFetchHtml
			return
		}

		if r.StatusCode != 200 {
			cs.logger.Errorw("Error fetching page", "status", r.StatusCode, "pageID", pageID)

			failedErr = domain.ErrCrawlFailedToFetchHtml
			return
		}

		pageHash = util.Sha256Hex32(r.Body)

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
			page, err = cs.pageService.Create(pageID, url, pageHash)
			if err != nil {
				cs.logger.Errorw("Error creating page", "error", err, "pageID", pageID)
				failedErr = err
			}
		}

		cs.logger.Infow("Page created", "pageID", page.ID, "url", page.URL)
		pageCrawled = page
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		dst := e.Attr("href")

		// if link is relative, make it absolute
		absoluteDst, err := util.MakeAbsoluteIfRelative(url, dst)

		if err != nil {
			cs.logger.Errorw("Error making link absolute", "error", err, "link", dst)
			return
		}

		if err := cs.api.CreateLink(url, absoluteDst); err != nil {
			cs.logger.Errorw("Error creating link", "error", err, "link", absoluteDst)
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		cs.logger.Errorw("Error crawling page", "error", e, "pageID", pageID)
		failedErr = e
	})

	c.Visit(url)

	return pageCrawled, failedErr
}

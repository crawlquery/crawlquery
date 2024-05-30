package service

import (
	"crawlquery/node/domain"
	"crawlquery/pkg/client/api"
	"crawlquery/pkg/util"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
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

func (cs *CrawlService) Crawl(pageID, url string) (contentHash string, links []string, failedErr error) {

	// Instantiate default collector
	c := colly.NewCollector()

	c.IgnoreRobotsTxt = false

	extensions.RandomUserAgent(c)

	c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		// Returning an error prevents the redirect
		return http.ErrUseLastResponse
	})

	c.OnResponse(func(r *colly.Response) {

		if !strings.Contains(r.Headers.Get("Content-Type"), "text/html") {
			cs.logger.Errorw("Error fetching page", "error", "Content-Type is not text/html", "pageID", pageID, "url", url)
			failedErr = domain.ErrCrawlFailedToFetchHtml
			return
		}

		if r.StatusCode != 200 {
			cs.logger.Errorw("Error fetching page", "status", r.StatusCode, "pageID", pageID, "url", url)

			failedErr = domain.ErrCrawlFailedToFetchHtml
			return
		}

		contentHash = util.Sha256Hex32(r.Body)

		err := cs.htmlService.Save(contentHash, r.Body)
		if err != nil {
			cs.logger.Errorw("Error saving page", "error", err, "pageID", pageID, "url", url)
			failedErr = domain.ErrCrawlFailedToStoreHtml
			return
		}

		cs.logger.Infow("Page crawled", "pageID", pageID, "url", url)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		dst := e.Attr("href")

		// if link is relative, make it absolute
		absoluteDst, err := util.MakeAbsoluteIfRelative(url, dst)

		if err != nil {
			cs.logger.Errorw("Error making link absolute", "error", err, "link", dst)
			return
		}

		links = append(links, absoluteDst)
	})

	c.OnError(func(r *colly.Response, e error) {
		cs.logger.Errorw("Error crawling page", "error", e, "pageID", pageID)
		failedErr = e
	})

	err := c.Visit(url)

	if err != nil {
		cs.logger.Errorw("Error visiting page", "error", err, "pageID", pageID)
		return "", nil, err
	}

	return
}

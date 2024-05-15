package service

import (
	"crawlquery/node/html/repository/disk"
	"crawlquery/node/index"
	"crawlquery/node/parse"
	"crawlquery/pkg/domain"
	"fmt"
	"io"
	"net/http"

	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

type CrawlService struct {
	idx            *index.Index
	htmlRepository *disk.Repository
	logger         *zap.SugaredLogger
}

func NewService(idx *index.Index, logger *zap.SugaredLogger) *CrawlService {
	return &CrawlService{
		idx:    idx,
		logger: logger,
	}
}

func (cs *CrawlService) Crawl(page *domain.Page) error {
	res, err := http.Get(page.URL)

	if err != nil {
		cs.logger.Errorw("Error fetching URL", "error", err, "url", url)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		cs.logger.Errorw("Error fetching URL", "status", res.StatusCode, "url", url)
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	html, err := io.ReadAll(res.Body)
	page, err := parse.Parse(html, url)

	// Instantiate default collector
	c := colly.NewCollector(
		// Attach a debugger to the collector
		colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnResponse(func(r *colly.Response) {
		
		cs.htmlRepository.Save(r.Body

		if err != nil {
			cs.logger.Errorw("Error reading response body", "error", err)
			return
		}

	})

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		//Delay:      5 * time.Second,
	})

	// Start scraping in five threads on https://httpbin.org/delay/2
	for i := 0; i < 5; i++ {
		c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	// Wait until threads are finished
	c.Wait()

	return nil
}

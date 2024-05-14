package service

import (
	"crawlquery/node/index"
	"crawlquery/node/parse"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type CrawlService struct {
	idx    *index.Index
	logger *zap.SugaredLogger
}

func NewService(idx *index.Index, logger *zap.SugaredLogger) *CrawlService {
	return &CrawlService{
		idx:    idx,
		logger: logger,
	}
}

func (cs *CrawlService) Crawl(url string) error {
	res, err := http.Get(url)

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
	return cs.idx.AddPage(page)
}

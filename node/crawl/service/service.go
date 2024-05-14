package service

import (
	"crawlquery/node/index"
	"crawlquery/pkg/domain"
	"crawlquery/pkg/util"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		cs.logger.Errorw("Error parsing HTML", "error", err, "url", url)
		return err
	}
	html, err := doc.Html()

	if err != nil {

		cs.logger.Errorw("Error getting HTML", "error", err, "url", url)
		return err
	}

	page := &domain.Page{
		ID:  util.UUID(),
		URL: url,
		// get the title from the head of the HTML document
		Title:           doc.Find("head title").Text(),
		Content:         html,
		MetaDescription: doc.Find("meta[name=description]").AttrOr("content", ""),
	}
	return cs.idx.AddPage(page)
}

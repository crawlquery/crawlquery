package service

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/util"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type CrawlService struct {
	indexService domain.IndexService
}

func NewCrawlService(is domain.IndexService) *CrawlService {
	return &CrawlService{
		indexService: is,
	}
}

func (cs *CrawlService) Crawl(url string) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return err
	}
	html, err := doc.Html()

	if err != nil {
		return err
	}

	page := &domain.Page{
		ID:              util.UUID(),
		URL:             url,
		Title:           doc.Find("title").Text(),
		Content:         html,
		MetaDescription: doc.Find("meta[name=description]").AttrOr("content", ""),
	}
	return cs.indexService.AddPage(page)
}

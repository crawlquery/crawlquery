package parse

import (
	"bytes"
	"crawlquery/node/token"
	"crawlquery/pkg/domain"

	"github.com/PuerkitoBio/goquery"
)

func Title(doc *goquery.Document) string {
	return doc.Find("head title").Text()
}

func MetaDescription(doc *goquery.Document) string {
	return doc.Find("head meta[name=description]").AttrOr("content", "")
}

func Parse(html []byte, url string) (*domain.Page, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	if err != nil {
		return nil, err
	}

	page := &domain.Page{
		URL: url,
	}

	page.Title = Title(doc)
	page.MetaDescription = MetaDescription(doc)
	page.Keywords = token.Keywords(doc)

	return page, nil
}

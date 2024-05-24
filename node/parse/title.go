package parse

import (
	"crawlquery/node/domain"

	"github.com/PuerkitoBio/goquery"
	"github.com/jpillora/go-tld"
)

type TitleParser struct {
	doc *goquery.Document
}

func NewTitleParser(doc *goquery.Document) *TitleParser {
	return &TitleParser{
		doc: doc,
	}
}

func (tp *TitleParser) Parse(page *domain.Page) {

	ogTitle := tp.doc.Find("meta[property='og:title']").AttrOr("content", "")

	if ogTitle != "" {
		page.Title = ogTitle
		return
	}

	titleTag := tp.doc.Find("title").Text()

	if titleTag != "" {
		page.Title = titleTag
		return
	}

	domain, err := tld.Parse(page.URL)

	if err != nil {
		page.Title = "We couldn't find a title for this page."
		return
	}

	page.Title = domain.Host
}

package parse

import (
	"crawlquery/node/domain"

	"github.com/PuerkitoBio/goquery"
)

type DescriptionParser struct {
	doc *goquery.Document
}

func NewDescriptionParser(doc *goquery.Document) *DescriptionParser {
	return &DescriptionParser{
		doc: doc,
	}
}

func (dp *DescriptionParser) Parse(page *domain.Page) {

	ogDescription := dp.doc.Find("meta[property='og:description']").AttrOr("content", "")

	if ogDescription != "" {
		page.Description = ogDescription
		return
	}

	metaDescription := dp.doc.Find("meta[name='description']").AttrOr("content", "")

	if metaDescription != "" {
		page.Description = metaDescription
		return
	}

	// first paragraph
	firstParagraph := dp.doc.Find("p").First().Text()

	if firstParagraph != "" {
		page.Description = firstParagraph
		return
	}

	page.Description = "We couldn't find a description for this page."
}

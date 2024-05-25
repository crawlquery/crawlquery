package parse

import (
	"crawlquery/node/domain"
	"errors"
	"strings"

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

func (dp *DescriptionParser) Parse(page *domain.Page) error {

	ogDescription := dp.doc.Find("meta[property='og:description']").AttrOr("content", "")

	if ogDescription != "" {
		page.Description = ogDescription
		return nil
	}

	metaDescription := dp.doc.Find("meta[name='description']").AttrOr("content", "")

	if metaDescription != "" {
		page.Description = metaDescription
		return nil
	}

	// first paragraph
	firstParagraph := dp.doc.Find("p").First().Text()

	// check string isn't just whitespace
	if strings.TrimSpace(firstParagraph) == "" {
		return errors.New("no description found")
	}

	if firstParagraph != "" {
		page.Description = firstParagraph
		return nil
	}

	return errors.New("no description found")
}

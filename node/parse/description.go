package parse

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Description(doc *goquery.Document) (string, error) {

	ogDescription := doc.Find("meta[property='og:description']").AttrOr("content", "")

	if ogDescription != "" {
		return ogDescription, nil
	}

	metaDescription := doc.Find("meta[name='description']").AttrOr("content", "")

	if metaDescription != "" {
		return metaDescription, nil
	}

	// first paragraph
	firstParagraph := doc.Find("p").First().Text()

	// check string isn't just whitespace
	if strings.TrimSpace(firstParagraph) == "" {
		firstParagraph = ""
	}

	if firstParagraph != "" {
		return firstParagraph, nil
	}

	return "", errors.New("no description found")
}

package parse

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

func Title(doc *goquery.Document) (string, error) {

	ogTitle := doc.Find("meta[property='og:title']").AttrOr("content", "")

	if ogTitle != "" {
		return ogTitle, nil
	}

	titleTag := doc.Find("title").Text()

	if titleTag != "" {
		return titleTag, nil
	}

	return "", errors.New("no title found")
}

package parse

import (
	"crawlquery/node/token"

	"github.com/PuerkitoBio/goquery"
)

func Title(doc *goquery.Document) string {
	return doc.Find("head title").Text()
}

func MetaDescription(doc *goquery.Document) string {
	return doc.Find("head meta[name=description]").AttrOr("content", "")
}

func Keywords(doc *goquery.Document) []string {
	return token.Keywords(doc)
}

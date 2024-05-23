package parse

import (
	"github.com/PuerkitoBio/goquery"
)

func Title(doc *goquery.Document) string {
	return doc.Find("head title").Text()
}

func MetaDescription(doc *goquery.Document) string {
	return doc.Find("head meta[name=description]").AttrOr("content", "")
}

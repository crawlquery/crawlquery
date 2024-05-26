package parse

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
)

func Language(doc *goquery.Document) (string, bool) {
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		Build()

	var textBuilder strings.Builder

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		textBuilder.WriteString(s.Text())
	})

	doc.Find("h1, h2, h3, h4, h5").Each(func(i int, s *goquery.Selection) {
		textBuilder.WriteString(s.Text())
	})

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		textBuilder.WriteString(s.Text())
	})

	lang, reliable := detector.DetectLanguageOf(textBuilder.String())

	return lang.String(), reliable
}

package parse

import (
	"crawlquery/node/domain"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
)

type LanguageParser struct {
	doc *goquery.Document
}

func NewLanguageParser(doc *goquery.Document) *LanguageParser {
	return &LanguageParser{
		doc: doc,
	}
}

func (lp *LanguageParser) Parse(page *domain.Page) error {
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		Build()

	var textBuilder strings.Builder

	lp.doc.Find("p").Each(func(i int, s *goquery.Selection) {
		textBuilder.WriteString(s.Text())
	})

	lp.doc.Find("h1, h2, h3, h4, h5").Each(func(i int, s *goquery.Selection) {
		textBuilder.WriteString(s.Text())
	})

	lp.doc.Find("a").Each(func(i int, s *goquery.Selection) {
		textBuilder.WriteString(s.Text())
	})

	lang, _ := detector.DetectLanguageOf(textBuilder.String())

	page.Language = lang.String()

	return nil
}

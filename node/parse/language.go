package parse

import (
	"crawlquery/node/domain"

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

	lang, _ := detector.DetectLanguageOf(lp.doc.Text())

	page.Language = lang.String()

	return nil
}

package parse

import (
	"crawlquery/node/domain"

	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
)

type LanguageParser struct {
	page *domain.Page
}

func NewLanguageParser(page *domain.Page) *LanguageParser {
	return &LanguageParser{
		page: page,
	}
}

func (lp *LanguageParser) Parse(doc *goquery.Document) {
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		Build()

	lang, _ := detector.DetectLanguageOf(doc.Text())

	lp.page.Language = lang.String()
}

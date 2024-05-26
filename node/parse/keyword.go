package parse

import (
	"crawlquery/node/domain"
	"crawlquery/node/keyword"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
)

type KeywordParser struct {
	doc      *goquery.Document
	keywords *[][]string
}

func NewKeywordParser(doc *goquery.Document, keywords *[][]string) *KeywordParser {
	return &KeywordParser{
		doc:      doc,
		keywords: keywords,
	}
}

func (pp *KeywordParser) HeadingKeywords() ([][]string, error) {
	var headings []string
	pp.doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		headings = append(headings, s.Text())
	})

	var parsedKeywords [][]string

	for _, h := range headings {
		clean := strings.Join(strings.Fields(h), " ")

		parsed, err := keyword.ParseText(clean)
		if err != nil {
			return nil, err
		}
		parsedKeywords = append(parsedKeywords, parsed...)
	}

	return parsedKeywords, nil
}

func (pp *KeywordParser) ParseParagraph() ([][]string, error) {

	var paragraphs []string
	pp.doc.Find("p").Each(func(i int, s *goquery.Selection) {
		paragraphs = append(paragraphs, s.Text())
	})

	var keywords [][]string

	for _, p := range paragraphs {
		clean := strings.Join(strings.Fields(p), " ")
		parsedKeywords, err := keyword.ParseText(clean)
		if err != nil {
			return nil, err
		}

		keywords = append(keywords, parsedKeywords...)
	}

	return keywords, nil
}

func (kp *KeywordParser) Parse(page *domain.Page) error {

	if page.Language != lingua.English.String() {
		return errors.New("only english pages are supported")
	}

	paragraphKeywords, err := kp.ParseParagraph()
	if err != nil {
		return err
	}

	headingKeywords, err := kp.HeadingKeywords()
	if err != nil {
		return err
	}

	var keywords [][]string

	keywords = append(keywords, paragraphKeywords...)
	keywords = append(keywords, headingKeywords...)

	*kp.keywords = keywords

	return nil
}

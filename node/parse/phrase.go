package parse

import (
	"crawlquery/node/domain"
	"crawlquery/node/phrase"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
)

type PhraseParser struct {
	doc *goquery.Document
}

func NewPhraseParser(doc *goquery.Document) *PhraseParser {
	return &PhraseParser{
		doc: doc,
	}
}

func (pp *PhraseParser) HeadingPhrases() ([][]string, error) {
	var headings []string
	pp.doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		headings = append(headings, s.Text())
	})

	var parsedPhrases [][]string

	for _, h := range headings {
		clean := strings.Join(strings.Fields(h), " ")

		parsed, err := phrase.ParseText(clean)
		if err != nil {
			return nil, err
		}
		parsedPhrases = append(parsedPhrases, parsed...)
	}

	return parsedPhrases, nil
}

func (pp *PhraseParser) ParseParagraph() ([][]string, error) {

	var paragraphs []string
	pp.doc.Find("p").Each(func(i int, s *goquery.Selection) {
		paragraphs = append(paragraphs, s.Text())
	})

	var phrases [][]string

	for _, p := range paragraphs {
		clean := strings.Join(strings.Fields(p), " ")
		parsedPhrases, err := phrase.ParseText(clean)
		if err != nil {
			return nil, err
		}

		phrases = append(phrases, parsedPhrases...)
	}

	return phrases, nil
}

func (kp *PhraseParser) Parse(page *domain.Page) error {

	if page.Language != lingua.English.String() {
		return errors.New("only english pages are supported")
	}

	paragraphPhrases, err := kp.ParseParagraph()
	if err != nil {
		return err
	}

	headingPhrases, err := kp.HeadingPhrases()
	if err != nil {
		return err
	}

	page.Phrases = append(page.Phrases, paragraphPhrases...)
	page.Phrases = append(page.Phrases, headingPhrases...)

	return nil
}

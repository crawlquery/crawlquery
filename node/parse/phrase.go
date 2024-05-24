package parse

import (
	"crawlquery/node/domain"
	"crawlquery/node/phrase"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/neurosnap/sentences/english"
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

func (pp *PhraseParser) HeadingPhrases() []string {
	var phrases []string

	pp.doc.Find("h1 h2 h3").Each(func(i int, s *goquery.Selection) {
		phrases = append(phrases, s.Text())
	})

	return phrases
}

func (pp *PhraseParser) ParseParagraph() ([][]string, error) {

	var paragraphs []string
	pp.doc.Find("p").Each(func(i int, s *goquery.Selection) {
		paragraphs = append(paragraphs, s.Text())
	})

	var phrases [][]string

	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return nil, err
	}

	for _, p := range paragraphs {
		sentences := tokenizer.Tokenize(p)
		for _, s := range sentences {
			parsedPhrases, err := phrase.ParseSentence(s.Text)

			if err != nil {
				return nil, err
			}

			phrases = append(phrases, parsedPhrases...)
		}
	}

	return phrases, nil
}

func (kp *PhraseParser) Parse(page *domain.Page) error {

	if page.Language != lingua.English.String() {
		return errors.New("only english pages are supported")
	}

	var phrases [][]string

	paragraphPhrases, err := kp.ParseParagraph()
	if err != nil {
		return err
	}

	page.Phrases = append(phrases, paragraphPhrases...)

	return nil
}

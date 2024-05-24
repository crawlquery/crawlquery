package parse

import (
	"crawlquery/node/domain"

	"github.com/PuerkitoBio/goquery"
	"github.com/jdkato/prose/v2"
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

var PhraseTemplates = [][]string{
	// the best five search engines
	{"DT", "JJS", "CD", "NN", "NNS"},
	// the best search engine
	{"DT", "JJS", "NN", "NN"},
}

func (PhraseParser) ParseSentence(sentence string) ([]string, error) {
	doc, err := prose.NewDocument(sentence)
	if err != nil {
		return nil, err
	}

	var phrase []string

	for i := 0; i < len(doc.Tokens()); i++ {
		for _, template := range PhraseTemplates {
			if i+len(template) > len(doc.Tokens()) {
				continue
			}

			match := true
			for j, pos := range template {
				if doc.Tokens()[i+j].Tag != pos {
					match = false
					break
				}
			}

			if match {
				for j := range template {
					phrase = append(phrase, doc.Tokens()[i+j].Text)
				}
				return phrase, nil
			}
		}
	}

	return phrase, nil
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
			parsed, err := pp.ParseSentence(s.Text)

			if err != nil {
				return nil, err
			}
			phrases = append(phrases, parsed)
		}
	}

	return phrases, nil
}

func (kp *PhraseParser) Parse(page *domain.Page) {

	if page.Language != lingua.EN.String() {
		return
	}
}

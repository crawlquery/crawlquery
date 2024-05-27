package parse

import (
	"crawlquery/node/domain"
	"crawlquery/node/keyword"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HeadingKeywords(doc *goquery.Document) ([]domain.Keyword, error) {
	var headings []string
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		headings = append(headings, s.Text())
	})

	var parsedKeywords []domain.Keyword

	for _, h := range headings {
		clean := strings.ToLower(strings.Join(strings.Fields(h), " "))

		parsed, err := keyword.ParseText(clean)
		if err != nil {
			return nil, err
		}
		parsedKeywords = append(parsedKeywords, parsed...)
	}

	return parsedKeywords, nil
}

func ParseParagraph(doc *goquery.Document) ([]domain.Keyword, error) {

	var paragraphs []string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		paragraphs = append(paragraphs, s.Text())
	})

	var keywords []domain.Keyword

	clean := strings.ToLower(strings.Join(paragraphs, " "))
	parsedKeywords, err := keyword.ParseText(clean)
	if err != nil {
		return nil, err
	}

	keywords = append(keywords, parsedKeywords...)

	return keywords, nil
}

func Keywords(doc *goquery.Document) ([]domain.Keyword, error) {

	paragraphKeywords, err := ParseParagraph(doc)
	if err != nil {
		return nil, err
	}

	headingKeywords, err := HeadingKeywords(doc)
	if err != nil {
		return nil, err
	}

	var keywords []domain.Keyword

	keywords = append(keywords, paragraphKeywords...)
	keywords = append(keywords, headingKeywords...)

	return keywords, nil
}

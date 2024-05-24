package parse_test

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"slices"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestPhraseParserParseSentence(t *testing.T) {
	t.Run("parses phrases from a sentence", func(t *testing.T) {
		cases := []struct {
			sentence string
			want     []string
		}{

			// Infinitive phrases
			{
				sentence: "What are the best five search engines, what are your dreams about?",
				want:     []string{"the", "best", "five", "search", "engines"},
			},
			{
				sentence: "What are the best search engine, what are your dreams about?",
				want:     []string{"the", "best", "search", "engine"},
			},
			{
				sentence: "What are the best search engines, what are your dreams about?",
				want:     []string{"the", "best", "search", "engines"},
			},

			// Noun phrases
			{
				sentence: "What are you up to today, do you want to know how to make a cake?",
				want:     []string{"how", "to", "make", "a", "cake"},
			},
		}

		for _, tc := range cases {
			t.Run("parses phrases", func(t *testing.T) {
				kp := parse.PhraseParser{}

				got, err := kp.ParseSentence(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				if !slices.Equal(got, tc.want) {
					t.Errorf("Expected %v, got %v", tc.want, got)
				}
			})
		}
	})
}

func TestPhraseParagraphPhrases(t *testing.T) {
	t.Run("parses phrases from a paragraph", func(t *testing.T) {
		cases := []struct {
			html []byte
			want [][]string
		}{
			{
				html: testdataloader.GetTestFile("testdata/pages/info/which-search-engine-is-the-best.html"),
				want: [][]string{
					{"the", "best", "five", "search", "engines"},
				},
			},
		}

		for _, tc := range cases {
			t.Run("parses phrases", func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
				if err != nil {
					t.Errorf("Error parsing html: %v", err)
				}

				kp := parse.NewPhraseParser(doc)

				page := &domain.Page{
					Language: lingua.EN.String(),
				}
				kp.Parse(page)

			})
		}
	})
}

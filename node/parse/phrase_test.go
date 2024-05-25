package parse_test

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestPhraseParser(t *testing.T) {
	t.Run("only parses english pages", func(t *testing.T) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(testdataloader.GetTestFile("testdata/pages/recipe/how-to-make-bolognese-sauce.html")))
		if err != nil {
			t.Fatalf("Error loading document: %v", err)
		}

		pp := parse.NewPhraseParser(doc)

		page := &domain.Page{
			URL:      "http://example.com",
			Language: lingua.Italian.String(),
		}

		pp.Parse(page)

		if len(page.Phrases) != 0 {
			t.Errorf("Expected no phrases, got %v", page.Phrases)
		}
	})

	t.Run("parses phrases from the page", func(t *testing.T) {
		cases := []struct {
			html     []byte
			contains [][]string
		}{
			{
				html: testdataloader.GetTestFile("testdata/pages/info/which-search-engine-is-the-best.html"),
				contains: [][]string{
					{"search", "engine"},
				},
			},
		}

		for _, tc := range cases {
			t.Run("parses phrases from the page", func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
				if err != nil {
					t.Fatalf("Error loading document: %v", err)
				}

				pp := parse.NewPhraseParser(doc)

				page := &domain.Page{
					URL:      "http://example.com",
					Language: lingua.English.String(),
				}

				err = pp.Parse(page)

				if err != nil {
					t.Fatalf("Error parsing phrases: %v", err)
				}

				for _, c := range tc.contains {
					found := false
					for _, p := range page.Phrases {
						if len(p) != len(c) {
							continue
						}
						for i, w := range c {
							if p[i] != w {
								break
							}
							if i == len(c)-1 {
								found = true
							}
						}
					}

					if !found {
						t.Errorf("Expected to find %v in %v", c, page.Phrases)
					}
				}
			})
		}
	})

	t.Run("parses phrases from the heading", func(t *testing.T) {
		cases := []struct {
			html     []byte
			contains [][]string
		}{
			{
				html: testdataloader.GetTestFile("testdata/pages/stackoverflow/best-way-to-detect-bot-from-user-agent.html"),
				contains: [][]string{
					{"Best"},
					{"way"},
					{"detect"},
					{"bot"},
					{"from"},
					{"Best", "way", "to", "detect", "bot", "from", "user", "agent"},
					{"user", "agent"},
				},
			},
		}

		for _, tc := range cases {
			t.Run("parses phrases from the page", func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
				if err != nil {
					t.Fatalf("Error loading document: %v", err)
				}

				pp := parse.NewPhraseParser(doc)

				page := &domain.Page{
					URL:      "http://example.com",
					Language: lingua.English.String(),
				}

				err = pp.Parse(page)

				if err != nil {
					t.Fatalf("Error parsing phrases: %v", err)
				}

				for _, c := range tc.contains {
					found := false
					for _, p := range page.Phrases {
						if len(p) != len(c) {
							continue
						}
						for i, w := range c {
							if p[i] != w {
								break
							}
							if i == len(c)-1 {
								found = true
							}
						}
					}

					if !found {
						t.Errorf("Expected to find %v in %v", c, page.Phrases)
					}
				}
			})
		}
	})
}

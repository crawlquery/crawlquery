package parse_test

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"testing"

	"github.com/PuerkitoBio/goquery"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestKeywordParser(t *testing.T) {
	t.Run("parses keywords from the page", func(t *testing.T) {
		cases := []struct {
			html     []byte
			contains []domain.Keyword
		}{
			{
				html: testdataloader.GetTestFile("testdata/pages/info/which-search-engine-is-the-best.html"),
				contains: []domain.Keyword{
					"search engine",
				},
			},
		}

		for _, tc := range cases {
			t.Run("parses keywords from the page", func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
				if err != nil {
					t.Fatalf("Error loading document: %v", err)
				}

				keywords, err := parse.Keywords(doc)

				if err != nil {
					t.Fatalf("Error parsing keywords: %v", err)
				}

				for _, c := range tc.contains {
					found := false
					for _, p := range keywords {
						if len(p) != len(c) {
							continue
						}
						for i, w := range c {
							if p[i] != byte(w) {
								break
							}
							if i == len(c)-1 {
								found = true
							}
						}
					}

					if !found {
						t.Errorf("Expected to find %v in %v", c, keywords)
					}
				}
			})
		}
	})

	t.Run("parses keywords from the heading", func(t *testing.T) {
		cases := []struct {
			html     []byte
			contains []domain.Keyword
		}{
			{
				html: testdataloader.GetTestFile("testdata/pages/stackoverflow/best-way-to-detect-bot-from-user-agent.html"),
				contains: []domain.Keyword{
					"best",
					"way",
					"detect",
					"bot",
					"from",
					"best way to detect bot from user agent",
					"user agent",
				},
			},
		}

		for _, tc := range cases {
			t.Run("parses keywords from the page", func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
				if err != nil {
					t.Fatalf("Error loading document: %v", err)
				}

				keywords, err := parse.Keywords(doc)

				if err != nil {
					t.Fatalf("Error parsing keywords: %v", err)
				}

				if err != nil {
					t.Fatalf("Error parsing keywords: %v", err)
				}

				for _, c := range tc.contains {
					found := false
					for _, p := range keywords {
						if len(p) != len(c) {
							continue
						}
						for i, w := range c {
							if p[i] != byte(w) {
								break
							}
							if i == len(c)-1 {
								found = true
							}
						}
					}

					if !found {
						t.Errorf("Expected to find %v in %v", c, keywords)
					}
				}
			})
		}
	})
}

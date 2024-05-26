package parse_test

import (
	"bytes"
	"crawlquery/node/parse"
	"testing"

	"github.com/PuerkitoBio/goquery"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestDescriptionParser(t *testing.T) {
	t.Run("parses the description from the meta tag", func(t *testing.T) {
		cases := []struct {
			html []byte
			want string
		}{
			{
				html: testdataloader.GetTestFile("testdata/pages/info/which-search-engine-is-the-best.html"),
				want: "Find out which search engines are effective.",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/info/what-are-some-types-of-search-engines.html"),
				want: "The basic types of search engines include: Web crawlers, meta, directories and hybrids. Within these basic types, there are many different methods used to retrieve information. Some common search engines include Google, Bing and Yahoo.",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/info/how-to-change-the-default-search-engine-on-all-browsers.html"),
				want: "Learning how to change the search engine default settings in most major web browsers can be a boon for productivity. We show you how, here.",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/info/ways-to-reuse-egg-cartons.html"),
				want: "Egg cartons, whether cardboard, plastic or Styrofoam, can be reused for lots of projects in your home or garden. Here are 10 eggs-citing ones to try.",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/dummy/paragraph-only.html"),
				want: "This is how you create a webpage without a description, containing only a paragraph.",
			},
		}

		for _, tc := range cases {
			t.Run(tc.want, func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
				if err != nil {
					t.Errorf("Error parsing html: %v", err)
				}

				description, err := parse.Description(doc)

				if err != nil {
					t.Errorf("Error parsing description: %v", err)
				}

				if description != tc.want {
					t.Errorf("Expected %s, got %s", tc.want, description)
				}
			})
		}
	})
}

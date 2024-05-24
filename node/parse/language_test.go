package parse_test

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"testing"

	"github.com/PuerkitoBio/goquery"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestLanguageParser(t *testing.T) {
	t.Run("parses the language of a page", func(t *testing.T) {
		cases := []struct {
			html []byte
			want string
		}{
			{
				html: testdataloader.GetTestFile("testdata/pages/language/english.html"),
				want: "English",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/language/spanish.html"),
				want: "Spanish",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/language/german.html"),
				want: "German",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/language/chinese.html"),
				want: "Chinese",
			},
			{
				html: testdataloader.GetTestFile("testdata/pages/language/french.html"),
				want: "French",
			},
		}

		for _, tc := range cases {
			t.Run(tc.want, func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
				if err != nil {
					t.Errorf("Error parsing html: %v", err)
				}
				lp := parse.NewLanguageParser(doc)
				page := &domain.Page{}
				lp.Parse(page)

				if page.Language != tc.want {
					t.Errorf("Expected %s, got %s", tc.want, page.Language)
				}
			})
		}
	})

}

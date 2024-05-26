package parse_test

import (
	"bytes"
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

				got, reliable := parse.Language(doc)

				if !reliable {
					t.Errorf("Expected reliable language detection")
				}

				if got != tc.want {
					t.Errorf("Expected %v, got %v", tc.want, got)
				}

			})
		}
	})

}

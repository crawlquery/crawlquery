package parse_test

import (
	"bytes"
	"crawlquery/node/parse"
	"testing"

	"github.com/PuerkitoBio/goquery"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestTitleParser(t *testing.T) {

	cases := []struct {
		html []byte
		want string
	}{
		{
			html: testdataloader.GetTestFile("testdata/pages/info/which-search-engine-is-the-best.html"),
			want: "Which Search Engine Is the Best?",
		},
		{
			html: testdataloader.GetTestFile("testdata/pages/info/what-are-some-types-of-search-engines.html"),
			want: "What Are Some Types of Search Engines?",
		},
		{
			html: testdataloader.GetTestFile("testdata/pages/info/how-to-change-the-default-search-engine-on-all-browsers.html"),
			want: "How to Change the Default Search Engine on All Browsers and Devices",
		},
		{
			html: testdataloader.GetTestFile("testdata/pages/info/ways-to-reuse-egg-cartons.html"),
			want: "Get Cracking! 10 Ways to Reuse Egg Cartons",
		},
	}

	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(tc.html))
			if err != nil {
				t.Errorf("Error parsing html: %v", err)
			}

			title, err := parse.Title(doc)

			if err != nil {
				t.Errorf("Error parsing title: %v", err)
			}

			if title != tc.want {
				t.Errorf("Expected %s, got %s", tc.want, title)
			}
		})
	}
}

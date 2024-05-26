package keyword

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestAdjectiveKeywords(t *testing.T) {
	t.Run("parses adjective keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "RB JJ intensifier adjective",
				sentence: "It is very interesting when it comes to the topic of politics.",
				want:     [][]string{{"very", "interesting"}, {"interesting"}},
			},
			{
				name:     "JJ adjective",
				sentence: "The quick brown fox jumps over the lazy dog.",
				want:     [][]string{{"quick"}, {"lazy"}},
			},
			{
				name:     "JJS JJ adjective",
				sentence: "The best way to detect bot from user agent.",
				want:     [][]string{{"best"}, {"user"}},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parseKeywords(doc.Tokens(), KeywordCategories{
					"adjective": adjectiveKeywordSubCategories(),
				})
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				sortKeywords(tc.want)
				sortKeywords(got)

				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("Expected %v, got %v", tc.want, got)
				}
			})
		}
	})
}

package keyword

import (
	"reflect"
	"testing"

	"crawlquery/node/domain"

	"github.com/jdkato/prose/v2"
)

func TestParseAdverbialKeywords(t *testing.T) {
	t.Run("parses adverbial keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     []domain.Keyword
		}{
			{
				name:     "simple adverbial keyword",
				sentence: "He left just a few months ago.",
				want:     []domain.Keyword{"just a few months ago"},
			},
			{
				name:     "multiple adverbial keywords",
				sentence: "She arrived just a few months ago and left right after the meeting.",
				want:     []domain.Keyword{"just a few months ago", "right after the meeting"},
			},
			{
				sentence: "The price of eggs is finally falling but it was at an all-time high just a few months ago.",
				want:     []domain.Keyword{"just a few months ago"},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parseKeywords(doc.Tokens(), KeywordCategories{
					"adverbial": KeywordSubCategories{
						"adverbial": AdverbialKeywordTemplates,
					},
				})

				sortKeywords(tc.want)
				sortKeywords(got)

				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("Expected %v, got %v", tc.want, got)
				}
			})
		}
	})
}

package keyword

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestParsePrepositionalKeywords(t *testing.T) {
	t.Run("parses prepositional keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "simple prepositional keyword",
				sentence: "The price of eggs is rising.",
				want:     [][]string{{"price", "of", "eggs"}, {"of", "eggs"}},
			},
			{
				name:     "prepositional keyword with determiner",
				sentence: "The price of the eggs is rising.",
				want:     [][]string{{"price", "of", "the", "eggs"}, {"of"}},
			},
			{
				name:     "simple prepositional keyword",
				sentence: "The stock is at an all-time high.",
				want:     [][]string{{"at", "an", "all-time", "high"}},
			},
			{
				name:     "prepositional keyword with determiner",
				sentence: "Price of eggs is finally falling but it was at an all-time high.",
				want:     [][]string{{"Price", "of", "eggs"}, {"of", "eggs"}, {"at", "an", "all-time", "high"}},
			},
			{
				name:     "simple prepositional keyword",
				sentence: "The note was underneath the egg tray.",
				want:     [][]string{{"underneath", "the", "egg", "tray"}},
			},
			{
				name:     "multiple prepositional keywords",
				sentence: "He met her at the park and they walked with friends in London.",
				want: [][]string{
					{"at", "the", "park"},
					{"with", "friends"},
					{"in", "London"},
				},
			},
			{
				name:     "IN (preposition)",
				sentence: "Where are you from?",
				want: [][]string{
					{"from"},
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parseKeywords(doc.Tokens(), KeywordCategories{
					"prepositional": KeywordSubCategories{
						"prepositional": PrepositionalKeywordTemplates,
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

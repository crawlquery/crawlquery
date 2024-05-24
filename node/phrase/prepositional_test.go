package phrase

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestParsePrepositionalPhrases(t *testing.T) {
	t.Run("parses prepositional phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "simple prepositional phrase",
				sentence: "The price of eggs is rising.",
				want:     [][]string{{"price", "of", "eggs"}, {"of", "eggs"}},
			},
			{
				name:     "prepositional phrase with determiner",
				sentence: "The price of the eggs is rising.",
				want:     [][]string{{"price", "of", "the", "eggs"}},
			},
			{
				name:     "simple prepositional phrase",
				sentence: "The stock is at an all-time high.",
				want:     [][]string{{"at", "an", "all-time", "high"}},
			},
			{
				name:     "multiple prepositional phrases",
				sentence: "He met her at the park and they walked with friends in London.",
				want: [][]string{
					{"at", "the", "park"},
					{"with", "friends"},
					{"in", "London"},
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parsePhrases(doc.Tokens(), PhraseCategories{
					"prepositional": PhraseSubCategories{
						"prepositional": PrepositionalPhraseTemplates,
					},
				})

				sortPhrases(tc.want)
				sortPhrases(got)

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

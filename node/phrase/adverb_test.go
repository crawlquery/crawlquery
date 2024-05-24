package phrase

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestParseAdverbialPhrases(t *testing.T) {
	t.Run("parses adverbial phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "simple adverbial phrase",
				sentence: "He left just a few months ago.",
				want:     [][]string{{"just", "a", "few", "months", "ago"}},
			},
			{
				name:     "multiple adverbial phrases",
				sentence: "She arrived just a few months ago and left right after the meeting.",
				want: [][]string{
					{"just", "a", "few", "months", "ago"},
					{"right", "after", "the", "meeting"},
				},
			},
			{
				sentence: "The price of eggs is finally falling but it was at an all-time high just a few months ago.",
				want:     [][]string{{"just", "a", "few", "months", "ago"}},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parsePhrases(doc.Tokens(), PhraseCategories{
					"adverbial": PhraseSubCategories{
						"adverbial": AdverbialPhraseTemplates,
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

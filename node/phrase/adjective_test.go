package phrase

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestAdjectivePhrases(t *testing.T) {
	t.Run("parses adjective phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "RB JJ intensifier adjective",
				sentence: "It is very interesting when it comes to the topic of politics.",
				want:     [][]string{{"very", "interesting"}},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parsePhrases(doc.Tokens(), PhraseCategories{
					"adjective": adjectivePhraseSubCategories(),
				})
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				sortPhrases(tc.want)
				sortPhrases(got)

				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("Expected %v, got %v", tc.want, got)
				}
			})
		}
	})
}

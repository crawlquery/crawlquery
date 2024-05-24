package phrase

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestParseNounPhrases(t *testing.T) {
	t.Run("parses simple noun phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			// Simple noun phrases
			{
				name:     "simple noun phrase",
				sentence: "A tree",
				want:     [][]string{{"tree"}},
			},
			{
				name:     "Multiple simple noun phrases",
				sentence: "A tree in the forest",
				want:     [][]string{{"tree"}, {"forest"}},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)

				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parsePhrases(doc.Tokens(), PhraseCategories{
					"noun": PhraseSubCategories{
						"simple_noun": SimpleNounTemplates,
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

	t.Run("parses adjective noun phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			// Adjective noun phrases
			{
				name:     "Adjective noun phrase",
				sentence: "A fast car went by.",
				want:     [][]string{{"fast", "car"}},
			},
			{
				name:     "Multiple adjective noun phrases",
				sentence: "The quick brown fox jumps over the lazy dog.",
				want:     [][]string{{"quick", "brown", "fox"}, {"lazy", "dog"}},
			},
			{
				name:     "Multiple adjective noun phrases",
				sentence: "The bright red car flew over the magical rainbow.",
				want:     [][]string{{"bright", "red", "car"}, {"red", "car"}, {"magical", "rainbow"}},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)

				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parsePhrases(doc.Tokens(), PhraseCategories{
					"noun": PhraseSubCategories{
						"adjective_noun": AdjectiveNounTemplates,
					},
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

	t.Run("parses noun phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			// Noun phrases
			{
				name:     "Noun phrase",
				sentence: "I walked past a bright red car, and saw a lazy dog.",
				want:     [][]string{{"bright", "red", "car"}, {"lazy", "dog"}, {"red", "car"}, {"dog"}, {"car"}},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)

				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parsePhrases(doc.Tokens(), PhraseCategories{
					"noun": nounPhraseSubCategories(),
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

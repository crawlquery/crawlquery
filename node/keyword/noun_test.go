package keyword

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestParseNounKeywords(t *testing.T) {
	t.Run("parses simple noun keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			// Simple noun keywords
			{
				name:     "simple noun keyword",
				sentence: "A tree",
				want:     [][]string{{"tree"}},
			},
			{
				name:     "Multiple simple noun keywords",
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

				got, err := parseKeywords(doc.Tokens(), KeywordCategories{
					"noun": KeywordSubCategories{
						"simple_noun": SimpleNounTemplates,
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

	t.Run("parses adjective noun keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			// Adjective noun keywords
			{
				name:     "Adjective noun keyword",
				sentence: "A fast car went by.",
				want:     [][]string{{"fast", "car"}},
			},
			{
				name:     "Multiple adjective noun keywords",
				sentence: "The quick brown fox jumps over the lazy dog.",
				want:     [][]string{{"quick", "brown", "fox"}, {"lazy", "dog"}},
			},
			{
				name:     "Multiple adjective noun keywords",
				sentence: "The bright red car flew over the magical rainbow.",
				want:     [][]string{{"bright", "red", "car"}, {"red", "car"}, {"magical", "rainbow"}},
			},
			{
				name:     "Adjective noun keyword",
				sentence: "Best way to detect bot from user agent?",
				want: [][]string{
					{"Best", "way", "to", "detect", "bot", "from", "user", "agent"},
					{"user", "agent"},
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
					"noun": KeywordSubCategories{
						"adjective_noun": AdjectiveNounTemplates,
					},
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

	t.Run("parses noun keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			// Noun keywords
			{
				name:     "Noun keyword",
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

				got, err := parseKeywords(doc.Tokens(), KeywordCategories{
					"noun": nounKeywordSubCategories(),
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

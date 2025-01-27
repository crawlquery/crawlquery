package keyword

import (
	"reflect"
	"testing"

	"crawlquery/node/domain"

	"github.com/jdkato/prose/v2"
)

func TestParseNounKeywords(t *testing.T) {
	t.Run("parses simple noun keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     []domain.Keyword
		}{
			// Simple noun keywords
			{
				name:     "simple noun keyword",
				sentence: "A tree",
				want:     []domain.Keyword{"tree"},
			},
			{
				name:     "Multiple simple noun keywords",
				sentence: "A tree in the forest",
				want:     []domain.Keyword{"tree", "forest"},
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
			want     []domain.Keyword
		}{
			// Adjective noun keywords
			{
				name:     "Adjective noun keyword",
				sentence: "A fast car went by.",
				want:     []domain.Keyword{"fast car"},
			},
			{
				name:     "Multiple adjective noun keywords",
				sentence: "The quick brown fox jumps over the lazy dog.",
				want:     []domain.Keyword{"quick brown fox", "lazy dog"},
			},
			{
				name:     "Multiple adjective noun keywords",
				sentence: "The bright red car flew over the magical rainbow.",
				want:     []domain.Keyword{"bright red car", "magical rainbow", "red car"},
			},
			{
				name:     "Adjective noun keyword",
				sentence: "Best way to detect bot from user agent?",
				want:     []domain.Keyword{"Best way to detect bot from user agent", "user agent"},
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
					t.Errorf("Expected %+v, got %+v", tc.want, got)
				}
			})
		}
	})

	t.Run("parses noun keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     []domain.Keyword
		}{
			// Noun keywords
			{
				name:     "Noun keyword",
				sentence: "I walked past a bright red car, and saw a lazy dog.",
				want:     []domain.Keyword{"bright red car", "lazy dog", "red car", "dog", "car"},
			},
			{
				name:     "Noun keyword",
				sentence: "The Nasdaq closes at record high.",
				want:     []domain.Keyword{"Nasdaq", "record high", "Nasdaq closes", "high"},
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

	t.Run("parses noun verb keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     []domain.Keyword
		}{
			// Noun verb keywords
			{
				name:     "Nasdaq closes at record high.",
				sentence: "Nasdaq closes at record high.",
				want:     []domain.Keyword{"Nasdaq closes"},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence, prose.WithSegmentation(false), prose.WithExtraction(false))
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got := parseSubCategories(doc.Tokens(), KeywordSubCategories{
					"noun_verb": NounVerbTemplates,
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

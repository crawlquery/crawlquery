package keyword

import (
	"reflect"
	"testing"

	"crawlquery/node/domain"

	"github.com/jdkato/prose/v2"
)

func TestVerbKeywords(t *testing.T) {
	t.Run("parses verb keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     []domain.Keyword
		}{
			{
				name:     "VBD (verb, past tense)",
				sentence: "He walked quickly.",
				want:     []domain.Keyword{"walked"},
			},
			{
				name:     "VBG (verb, gerund or present participle)",
				sentence: "I enjoy walking.",
				want:     []domain.Keyword{"enjoy", "walking"},
			},
			{
				name:     "VBN (verb, past participle)",
				sentence: "He eaten quickly.",
				want:     []domain.Keyword{"eaten"},
			},
			{
				name:     "VBP (verb, non-3rd person singular present)",
				sentence: "They run every day.",
				want:     []domain.Keyword{"run"},
			},
			{
				name:     "VBZ (verb, 3rd person singular present)",
				sentence: "He runs quickly.",
				want:     []domain.Keyword{"runs"},
			},
			{
				name:     "VB (verb, base form)",
				sentence: "I want to detect this.",
				want:     []domain.Keyword{"want", "detect"},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got, err := parseKeywords(doc.Tokens(), KeywordCategories{
					"verb": verbKeywordSubCategories(),
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

	t.Run("parses verb keywords with adverbs", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     []domain.Keyword
		}{
			{
				name:     "VBD (verb, past tense)",
				sentence: "He walked quickly.",
				want:     []domain.Keyword{"walked quickly"},
			},
			{
				name:     "VBG (verb, gerund or present participle)",
				sentence: "I enjoy walking quickly.",
				want:     []domain.Keyword{"walking quickly"},
			},
			{
				name:     "VBN (verb, past participle)",
				sentence: "He eaten quickly.",
				want:     []domain.Keyword{"eaten quickly"},
			},
			{
				name:     "VBP (verb, non-3rd person singular present)",
				sentence: "They run quickly every day.",
				want:     []domain.Keyword{"run quickly"},
			},
			{
				name:     "VBZ (verb, 3rd person singular present)",
				sentence: "He runs quickly.",
				want:     []domain.Keyword{"runs quickly"},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				doc, err := prose.NewDocument(tc.sentence)
				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

				got := parseSubCategories(doc.Tokens(), KeywordSubCategories{
					"verb": VerbAdverbTemplates,
				})

				sortKeywords(tc.want)
				sortKeywords(got)

				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("Expected %v, got %v", tc.want, got)
				}
			})
		}
	})
}

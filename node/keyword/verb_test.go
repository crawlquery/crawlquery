package keyword

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestVerbKeywords(t *testing.T) {
	t.Run("parses verb keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "VBD (verb, past tense)",
				sentence: "He walked quickly.",
				want:     [][]string{{"walked"}},
			},
			{
				name:     "VBG (verb, gerund or present participle)",
				sentence: "I enjoy walking.",
				want:     [][]string{{"enjoy"}, {"walking"}},
			},
			{
				name:     "VBN (verb, past participle)",
				sentence: "He eaten quickly.",
				want:     [][]string{{"eaten"}},
			},
			{
				name:     "VBP (verb, non-3rd person singular present)",
				sentence: "They run every day.",
				want:     [][]string{{"run"}},
			},
			{
				name:     "VBZ (verb, 3rd person singular present)",
				sentence: "He runs quickly.",
				want:     [][]string{{"runs"}},
			},
			{
				name:     "VB (verb, base form)",
				sentence: "I want to detect this.",
				want:     [][]string{{"want"}, {"detect"}},
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
			want     [][]string
		}{
			{
				name:     "VBD (verb, past tense)",
				sentence: "He walked quickly.",
				want:     [][]string{{"walked", "quickly"}},
			},
			{
				name:     "VBG (verb, gerund or present participle)",
				sentence: "I enjoy walking quickly.",
				want:     [][]string{{"walking", "quickly"}},
			},
			{
				name:     "VBN (verb, past participle)",
				sentence: "He eaten quickly.",
				want:     [][]string{{"eaten", "quickly"}},
			},
			{
				name:     "VBP (verb, non-3rd person singular present)",
				sentence: "They run quickly every day.",
				want:     [][]string{{"run", "quickly"}},
			},
			{
				name:     "VBZ (verb, 3rd person singular present)",
				sentence: "He runs quickly.",
				want:     [][]string{{"runs", "quickly"}},
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

package phrase

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestWordClasses(t *testing.T) {
	sentence := "He met her at the park and they walked with friends."
	doc, _ := prose.NewDocument(sentence)

	tokens := doc.Tokens()

	fmt.Printf("Tokens: %v\n", tokens)

	t.Fail()
}

func TestParseSentence(t *testing.T) {
	t.Run("parses a sentenc using noun, verb, and adjective phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "simple sentence",
				sentence: "The price of eggs is finally falling but it was at an all-time high just a few months ago.",
				want:     [][]string{{"price", "of", "eggs"}, {"falling"}, {"is"}, {"just", "a", "few", "months", "ago"}, {"was"}, {"at", "an", "all-time", "high"}},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := ParseSentence(tc.sentence)
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

func sortPhrases(phrases [][]string) {
	sort.Slice(phrases, func(i, j int) bool {
		return phraseString(phrases[i]) < phraseString(phrases[j])
	})
}

func phraseString(phrase []string) string {
	return strings.Join(phrase, " ")
}

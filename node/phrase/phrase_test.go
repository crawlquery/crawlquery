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
	sentence := "First, add some drainage for water by poking a few holes in the bottom of the carton. Barton Hill Farms suggests separating the lid from the bottom, then putting it underneath the egg tray to catch any wayward water."
	doc, _ := prose.NewDocument(sentence)

	tokens := doc.Tokens()

	fmt.Printf("Tokens: %v\n", tokens)

	// t.Fail()
}

func TestParseSentence(t *testing.T) {
	t.Run("parses a sentence using noun, verb, and adjective phrases", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     [][]string
		}{
			{
				name:     "price of eggs",
				sentence: "The price of eggs is finally falling but it was at an all-time high just a few months ago.",
				want:     [][]string{{"price", "of", "eggs"}, {"of", "eggs"}, {"price"}, {"falling"}, {"is"}, {"just", "a", "few", "months", "ago"}, {"few", "months"}, {"was"}, {"at", "an", "all-time", "high"}},
			},
			{
				name:     "add some drainage",
				sentence: "First, add some drainage for water by poking a few holes in the bottom of the carton. Barton Hill Farms suggests separating the lid from the bottom, then putting it underneath the egg tray to catch any wayward water.",
				want: [][]string{
					{"bottom"},
					{"bottom"},
					{"carton"},
					{"drainage"},
					{"egg", "tray"},
					{"few", "holes"},
					{"from", "the", "bottom"},
					{"in", "the", "bottom"},
					{"lid"},
					{"of", "the", "carton"},
					{"poking"},
					{"putting"},
					{"separating"},
					{"suggests"},
					{"tray"},
					{"underneath", "the", "egg", "tray"},
					{"water"},
					{"water"},
					{"wayward", "water"},
				},
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

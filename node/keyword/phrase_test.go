package keyword

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestWordClasses(t *testing.T) {
	text := "Best way to detect bot from user agent?"
	doc, _ := prose.NewDocument(text)

	tokens := doc.Tokens()

	fmt.Printf("Tokens: %v\n", tokens)

	// t.Fail()
}

func TestParseText(t *testing.T) {
	t.Run("parses a text using noun, verb, and adjective keywords", func(t *testing.T) {
		cases := []struct {
			name string
			text string
			want [][]string
		}{
			{
				name: "price of eggs",
				text: "The price of eggs is finally falling but it was at an all-time high just a few months ago.",
				want: [][]string{{"price", "of", "eggs"}, {"of", "eggs"}, {"price"}, {"falling"}, {"is"}, {"just", "a", "few", "months", "ago"}, {"few", "months"}, {"was"}, {"at", "an", "all-time", "high"}},
			},
			{
				name: "add some drainage",
				text: "First, add some drainage for water by poking a few holes in the bottom of the carton. Barton Hill Farms suggests separating the lid from the bottom, then putting it underneath the egg tray to catch any wayward water.",
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
				got, err := ParseText(tc.text)
				if err != nil {
					t.Errorf("Error parsing text: %v", err)
				}

				sortKeywords(tc.want)
				sortKeywords(got)

				for _, keyword := range tc.want {
					var found bool

					for _, p := range got {
						if reflect.DeepEqual(keyword, p) {
							found = true
							break
						}
					}

					if !found {
						t.Errorf("Expected %v, got %v", tc.want, got)
					}
				}

			})
		}
	})

	t.Run("parses longest keyword matches", func(t *testing.T) {
		cases := []struct {
			name string
			text string
			want [][]string
		}{
			{
				name: "",
				text: "Best way to detect bot from user agent?",
				want: [][]string{
					{"Best", "way", "to", "detect", "bot", "from", "user", "agent"},
					{"agent"},
					{"bot"},
					{"user", "agent"},
					{"way"},
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := ParseText(tc.text)
				if err != nil {
					t.Errorf("Error parsing text: %v", err)
				}

				sortKeywords(tc.want)
				sortKeywords(got)

				for _, keyword := range tc.want {
					var found bool

					for _, p := range got {
						if reflect.DeepEqual(keyword, p) {
							found = true
							break
						}
					}

					if !found {
						t.Errorf("Expected %v, got %v", tc.want, got)
					}
				}
			})
		}
	})
}

func sortKeywords(keywords [][]string) {
	sort.Slice(keywords, func(i, j int) bool {
		return keywordString(keywords[i]) < keywordString(keywords[j])
	})
}

func keywordString(keyword []string) string {
	return strings.Join(keyword, " ")
}

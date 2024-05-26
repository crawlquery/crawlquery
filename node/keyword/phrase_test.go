package keyword

import (
	"crawlquery/node/domain"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestWordClasses(t *testing.T) {
	text := "Nasdaq closes at record high."
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
			want []domain.Keyword
		}{
			{
				name: "price of eggs",
				text: "The price of eggs is finally falling but it was at an all-time high just a few months ago.",
				want: []domain.Keyword{
					"price of eggs",
					"of eggs",
					"price",
					"falling",
					"is",
					"just a few months ago",
					"few months",
					"was",
					"at an all-time high",
					"high",
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
			want []domain.Keyword
		}{
			{
				name: "",
				text: "Best way to detect bot from user agent?",
				want: []domain.Keyword{
					"Best way to detect bot from user agent",
					"agent",
					"bot",
					"user", "agent",
					"way",
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

	t.Run("parses noun keywords", func(t *testing.T) {
		cases := []struct {
			name     string
			sentence string
			want     []domain.Keyword
		}{
			// Noun keywords
			{
				name:     "Noun keyword",
				sentence: "Nasdaq closes Friday at record high as Nvidia and the AI trade rallies on",
				want:     []domain.Keyword{"Nasdaq closes", "closes", "Friday", "Nvidia", "rallies"},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {

				got, err := ParseText(tc.sentence)

				if err != nil {
					t.Errorf("Error parsing sentence: %v", err)
				}

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

func sortKeywords(keywords []domain.Keyword) {
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i] < keywords[j]
	})
}

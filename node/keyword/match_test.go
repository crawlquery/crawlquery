package keyword

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestFindMatches(t *testing.T) {
	t.Run("finds matches", func(t *testing.T) {
		tokens := []prose.Token{
			{Tag: "DT", Text: "The"},
			{Tag: "JJ", Text: "quick"},
			{Tag: "NN", Text: "brown"},
			{Tag: "NN", Text: "fox"},
			{Tag: "VBZ", Text: "jumps"},
			{Tag: "IN", Text: "over"},
			{Tag: "DT", Text: "the"},
			{Tag: "JJ", Text: "lazy"},
			{Tag: "NN", Text: "dog"},
		}

		templates := KeywordSubCategory{
			{"DT", "JJ", "NN", "NN"},
			{"DT", "JJ", "NN"},
			{"DT", "NN"},
		}

		expected := []match{
			{start: 0, end: 3, keyword: []string{"The", "quick", "brown", "fox"}},
			{start: 0, end: 2, keyword: []string{"The", "quick", "brown"}},
			{start: 6, end: 8, keyword: []string{"the", "lazy", "dog"}},
		}

		got := findMatches(tokens, templates, 0)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	})

	t.Run("finds no matches if none exist", func(t *testing.T) {
		tokens := []prose.Token{
			{Tag: "VBZ", Text: "jumps"},
			{Tag: "IN", Text: "over"},
		}

		templates := KeywordSubCategory{
			{"DT", "JJ", "NN", "NN"},
			{"DT", "NN"},
		}

		expected := []match{}

		got := findMatches(tokens, templates, 0)

		if len(got) != 0 {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	})

	t.Run("finds matches starting from a specific index", func(t *testing.T) {
		tokens := []prose.Token{
			{Tag: "DT", Text: "The"},
			{Tag: "JJ", Text: "quick"},
			{Tag: "NN", Text: "brown"},
			{Tag: "NN", Text: "fox"},
			{Tag: "VBZ", Text: "jumps"},
			{Tag: "IN", Text: "over"},
			{Tag: "DT", Text: "the"},
			{Tag: "JJ", Text: "lazy"},
			{Tag: "NN", Text: "dog"},
		}

		templates := KeywordSubCategory{
			{"DT", "JJ", "NN", "NN"},
			{"DT", "JJ", "NN"},
			{"DT", "NN"},
		}

		expected := []match{
			{start: 6, end: 8, keyword: []string{"the", "lazy", "dog"}},
		}

		got := findMatches(tokens, templates, 6)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	})
}

func TestUpdateLongestMatches(t *testing.T) {
	t.Run("updates longest matches", func(t *testing.T) {
		longestMatches := map[int]match{
			0: {start: 0, end: 2, keyword: []string{"The", "quick"}},
		}

		matches := []match{
			{start: 0, end: 3, keyword: []string{"The", "quick", "brown", "fox"}},
			{start: 1, end: 3, keyword: []string{"quick", "brown", "fox"}},
			{start: 0, end: 2, keyword: []string{"The", "quick"}},
		}

		expected := map[int]match{
			0: {start: 0, end: 3, keyword: []string{"The", "quick", "brown", "fox"}},
			1: {start: 1, end: 3, keyword: []string{"quick", "brown", "fox"}},
		}

		updateLongestMatches(longestMatches, matches)

		if !reflect.DeepEqual(longestMatches, expected) {
			t.Errorf("Expected %v, got %v", expected, longestMatches)
		}
	})

	t.Run("adds new matches", func(t *testing.T) {
		longestMatches := map[int]match{}

		matches := []match{
			{start: 0, end: 2, keyword: []string{"The", "quick"}},
			{start: 1, end: 3, keyword: []string{"quick", "brown", "fox"}},
		}

		expected := map[int]match{
			0: {start: 0, end: 2, keyword: []string{"The", "quick"}},
			1: {start: 1, end: 3, keyword: []string{"quick", "brown", "fox"}},
		}

		updateLongestMatches(longestMatches, matches)

		if !reflect.DeepEqual(longestMatches, expected) {
			t.Errorf("Expected %v, got %v", expected, longestMatches)
		}
	})

	t.Run("does not update with shorter matches", func(t *testing.T) {
		longestMatches := map[int]match{
			0: {start: 0, end: 3, keyword: []string{"The", "quick", "brown", "fox"}},
		}

		matches := []match{
			{start: 0, end: 2, keyword: []string{"The", "quick"}},
		}

		expected := map[int]match{
			0: {start: 0, end: 3, keyword: []string{"The", "quick", "brown", "fox"}},
		}

		updateLongestMatches(longestMatches, matches)

		if !reflect.DeepEqual(longestMatches, expected) {
			t.Errorf("Expected %v, got %v", expected, longestMatches)
		}
	})
}

package phrase

import (
	"github.com/jdkato/prose/v2"
)

type CategoryName string
type SubCategoryName string
type Word string

type PhraseTemplate []Word
type PhraseSubCategory []PhraseTemplate
type PhraseCategories map[CategoryName]PhraseSubCategories
type PhraseSubCategories map[SubCategoryName]PhraseSubCategory

func phraseCategories() PhraseCategories {
	return PhraseCategories{
		"noun":          nounPhraseSubCategories(),
		"verb":          verbPhraseSubCategories(),
		"adjective":     adjectivePhraseSubCategories(),
		"prepositional": prepositionalPhraseSubCategories(),
		"adverbial":     adverbialPhraseSubCategories(),
		"quantifier":    quantifierPhraseSubCategories(),
	}
}

func ParseText(text string) ([][]string, error) {
	doc, err := prose.NewDocument(text, prose.WithSegmentation(false), prose.WithExtraction(false))
	if err != nil {
		return nil, err
	}
	tokens := doc.Tokens()

	return parsePhrases(tokens, phraseCategories())
}

type match struct {
	start  int
	end    int
	phrase []string
}

func parsePhrases(tokens []prose.Token, phraseCategories PhraseCategories) ([][]string, error) {

	var phrases [][]string

	for _, subCategories := range phraseCategories {
		subCategoryPhrases := parseSubCategories(tokens, subCategories)
		phrases = append(phrases, subCategoryPhrases...)
	}

	return phrases, nil
}

func parseSubCategories(tokens []prose.Token, subCategories PhraseSubCategories) [][]string {
	var phrases [][]string
	for _, subCategory := range subCategories {
		longestMatches := map[int]match{}

		for i := 0; i < len(tokens); i++ {
			matchedPhrases := findMatches(tokens, subCategory, i)
			updateLongestMatches(longestMatches, matchedPhrases)
		}

		// Convert the longest matches map to a slice
		for _, m := range longestMatches {
			phrases = append(phrases, m.phrase)
		}
	}
	return phrases
}

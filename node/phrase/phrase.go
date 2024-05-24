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

var SimpleNounTemplates = PhraseSubCategory{
	// A tree
	{"NN"},

	// A search engine
	{"NN", "NN"},
}

var AdjectiveNounTemplates = PhraseSubCategory{
	// bright red car
	{"JJ", "JJ", "NN"},

	// lazy dog
	{"JJ", "NN"},

	//quick brown fox
	{"JJ", "NN", "NN"},
}

func NounTemplates() PhraseSubCategories {
	templates := PhraseSubCategories{
		"simple_noun":    SimpleNounTemplates,
		"adjective_noun": AdjectiveNounTemplates,
	}

	return templates
}

func Templates() PhraseCategories {
	return PhraseCategories{
		"noun": NounTemplates(),
	}
}

type match struct {
	start  int
	end    int
	phrase []string
}

func parsePhrases(sentence string, phraseCategories PhraseCategories) ([][]string, error) {
	doc, err := prose.NewDocument(sentence)
	if err != nil {
		return nil, err
	}

	var phrases [][]string
	tokens := doc.Tokens()

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

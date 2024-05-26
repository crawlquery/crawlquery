package keyword

import (
	"crawlquery/node/domain"

	"github.com/jdkato/prose/v2"
)

type CategoryName string
type SubCategoryName string
type Word string

type KeywordTemplate []Word
type KeywordSubCategory []KeywordTemplate
type KeywordCategories map[CategoryName]KeywordSubCategories
type KeywordSubCategories map[SubCategoryName]KeywordSubCategory

func keywordCategories() KeywordCategories {
	return KeywordCategories{
		"noun":          nounKeywordSubCategories(),
		"verb":          verbKeywordSubCategories(),
		"adjective":     adjectiveKeywordSubCategories(),
		"prepositional": prepositionalKeywordSubCategories(),
		"adverbial":     adverbialKeywordSubCategories(),
		"quantifier":    quantifierKeywordSubCategories(),
	}
}

func ParseText(text string) ([]domain.Keyword, error) {
	doc, err := prose.NewDocument(text, prose.WithSegmentation(false), prose.WithExtraction(false))
	if err != nil {
		return nil, err
	}
	tokens := doc.Tokens()

	return parseKeywords(tokens, keywordCategories())
}

type match struct {
	start   int
	end     int
	keyword domain.Keyword
}

func parseKeywords(tokens []prose.Token, keywordCategories KeywordCategories) ([]domain.Keyword, error) {

	var keywords []domain.Keyword

	for _, subCategories := range keywordCategories {
		subCategoryKeywords := parseSubCategories(tokens, subCategories)
		keywords = append(keywords, subCategoryKeywords...)
	}

	return keywords, nil
}

func parseSubCategories(tokens []prose.Token, subCategories KeywordSubCategories) []domain.Keyword {
	var keywords []domain.Keyword
	for _, subCategory := range subCategories {
		longestMatches := map[int]match{}

		for i := 0; i < len(tokens); i++ {
			matchedKeywords := findMatches(tokens, subCategory, i)
			updateLongestMatches(longestMatches, matchedKeywords)
		}

		// Convert the longest matches map to a slice
		for _, m := range longestMatches {
			keywords = append(keywords, m.keyword)
		}
	}
	return keywords
}

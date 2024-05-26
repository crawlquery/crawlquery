package keyword

import (
	"crawlquery/node/domain"
	"strings"

	"github.com/jdkato/prose/v2"
)

func findMatches(tokens []prose.Token, templates KeywordSubCategory, startIndex int) []match {
	var matches []match
	for i := startIndex; i < len(tokens); i++ {
		for _, template := range templates {
			if i+len(template) <= len(tokens) {
				matchBool := true
				for j, wordType := range template {
					if tokens[i+j].Tag != string(wordType) {
						matchBool = false
						break
					}
				}
				if matchBool {
					var keywords []string
					for k := 0; k < len(template); k++ {
						keywords = append(keywords, tokens[i+k].Text)
					}

					kwString := strings.Join(keywords, " ")
					m := match{
						start:   i,
						end:     i + len(template) - 1,
						keyword: domain.Keyword(kwString),
					}
					matches = append(matches, m)
				}
			}
		}
	}
	return matches
}

func updateLongestMatches(longestMatches map[int]match, matches []match) {
	for _, m := range matches {
		if existingMatch, exists := longestMatches[m.start]; exists {
			if len(existingMatch.keyword) < len(m.keyword) {
				longestMatches[m.start] = m
			}
		} else {
			longestMatches[m.start] = m
		}
	}
}

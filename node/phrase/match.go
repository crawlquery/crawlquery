package phrase

import "github.com/jdkato/prose/v2"

func findMatches(tokens []prose.Token, templates PhraseSubCategory, startIndex int) []match {
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
					var phrase []string
					for k := 0; k < len(template); k++ {
						phrase = append(phrase, tokens[i+k].Text)
					}
					m := match{
						start:  i,
						end:    i + len(template) - 1,
						phrase: phrase,
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
			if len(existingMatch.phrase) < len(m.phrase) {
				longestMatches[m.start] = m
			}
		} else {
			longestMatches[m.start] = m
		}
	}
}

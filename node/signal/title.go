package signal

import (
	"crawlquery/node/domain"
	"regexp"
	"strings"
)

type Title struct{}

func (t *Title) anyMatch(title string, terms []string) domain.SignalLevel {
	baseLevel := domain.SignalLevelNone

	splitTitle := strings.Split(title, " ")
	for _, term := range terms {
		for _, titleWord := range splitTitle {
			if strings.EqualFold(titleWord, term) {
				baseLevel += domain.SignalLevelMedium
			}
		}
	}
	return baseLevel
}

func (t *Title) fullMatch(title string, terms []string) domain.SignalLevel {
	combined := strings.Join(terms, " ")
	if strings.EqualFold(title, combined) {
		return domain.SignalLevelMax
	}

	return domain.SignalLevelNone
}

func (ts *Title) Level(page *domain.Page, terms []string) domain.SignalLevel {

	baseLevel := domain.SignalLevelNone

	cleanedTitle := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(page.Title, "")

	cleanedTerms := make([]string, len(terms))

	for i, term := range terms {
		cleanedTerms[i] = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(term, "")
	}

	baseLevel += ts.anyMatch(cleanedTitle, cleanedTerms)

	baseLevel += ts.fullMatch(cleanedTitle, cleanedTerms)

	return baseLevel
}

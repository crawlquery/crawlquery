package signal

import "crawlquery/node/domain"

type Phrase struct{}

func (Phrase) Name() string {
	return "phrase"
}

func (p *Phrase) containsPhrase(page *domain.Page, phrase []string) bool {
	for _, pagePhrase := range page.Phrases {
		if len(pagePhrase) != len(phrase) {
			continue
		}

		matches := true

		for i, word := range pagePhrase {
			if word != phrase[i] {
				matches = false
				break
			}
		}

		if matches {
			return true
		}
	}

	return false
}

func (p *Phrase) Level(page *domain.Page, terms []string) (domain.SignalLevel, domain.SignalBreakdown) {

	baseLevel := domain.SignalLevelNone

	// create groups of terms starting from 1 to len(terms)
	// so for example if terms = ["a", "b", "c"]
	// groups = [["a"], ["a", "b"], ["a", "b", "c"],

	// then search the group for any matching phrases which could be ["a", "b"] or ["b", "c"]

	groups := make([][]string, len(terms)*(len(terms)+1)/2)

	for i := 0; i < len(terms); i++ {
		for j := i; j < len(terms); j++ {
			groups = append(groups, terms[i:j+1])
		}
	}

	for _, group := range groups {
		if p.containsPhrase(page, group) {
			baseLevel += domain.SignalLevelMedium
		}
	}

	return baseLevel, domain.SignalBreakdown{
		"phrase": baseLevel,
	}
}

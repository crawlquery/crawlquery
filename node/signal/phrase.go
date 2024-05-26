package signal

import "crawlquery/node/domain"

type Keyword struct{}

func (Keyword) Name() string {
	return "keyword"
}

func (p *Keyword) containsKeyword(page *domain.Page, keyword []string) bool {
	for _, pageKeyword := range page.Keywords {
		if len(pageKeyword) != len(keyword) {
			continue
		}

		matches := true

		for i, word := range pageKeyword {
			if word != keyword[i] {
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

func (p *Keyword) Level(page *domain.Page, terms []string) (domain.SignalLevel, domain.SignalBreakdown) {

	baseLevel := domain.SignalLevelNone

	// create groups of terms starting from 1 to len(terms)
	// so for example if terms = ["a", "b", "c"]
	// groups = [["a"], ["a", "b"], ["a", "b", "c"],

	// then search the group for any matching keywords which could be ["a", "b"] or ["b", "c"]

	groups := make([][]string, len(terms)*(len(terms)+1)/2)

	for i := 0; i < len(terms); i++ {
		for j := i; j < len(terms); j++ {
			groups = append(groups, terms[i:j+1])
		}
	}

	for _, group := range groups {
		if p.containsKeyword(page, group) {
			baseLevel += domain.SignalLevelMedium
		}
	}

	return baseLevel, domain.SignalBreakdown{
		"keyword": baseLevel,
	}
}

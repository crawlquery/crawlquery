package signal

import "crawlquery/node/domain"

type Keyword struct{}

func (ks *Keyword) Level(page *domain.Page, terms []string) domain.SignalLevel {
	baseLevel := domain.SignalLevelNone

	for _, term := range terms {
		for _, kw := range page.Keywords {
			if kw == term {
				baseLevel += domain.SignalLevelLow
			}
		}
	}

	return baseLevel
}

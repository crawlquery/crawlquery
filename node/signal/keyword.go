package signal

import "crawlquery/node/domain"

type Keyword struct{}

func (Keyword) Name() string {
	return "keyword"
}

func (ks *Keyword) Level(page *domain.Page, terms []string) (domain.SignalLevel, domain.SignalBreakdown) {
	baseLevel := domain.SignalLevelNone

	for _, term := range terms {
		for _, kw := range page.Keywords {
			if kw == term {
				baseLevel += domain.SignalLevelMedium
			}
		}
	}

	return baseLevel, domain.SignalBreakdown{
		"total": baseLevel,
	}
}

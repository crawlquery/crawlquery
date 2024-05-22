package signal

import "crawlquery/node/domain"

type KeywordSignal struct {
	keywordService domain.KeywordService
}

func (ks *KeywordSignal) Level(page string, keywords []string) domain.SignalLevel {
	return domain.SignalLevelNone
}

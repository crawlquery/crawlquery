package signal

import (
	"crawlquery/node/domain"

	tld "github.com/jpillora/go-tld"
)

// embed:

type Domain struct{}

func (ds *Domain) domainMatch(url *tld.URL, term []string) domain.SignalLevel {
	// domain matching
	for _, t := range term {
		if url.Domain == t {
			return domain.SignalLevelVeryHigh
		}
	}

	return domain.SignalLevelNone
}

func (ds *Domain) hostnameMatch(url *tld.URL, term []string) domain.SignalLevel {
	// hostname matching
	for _, t := range term {
		if url.Host == t {
			return domain.SignalLevelMax
		}
	}

	return domain.SignalLevelNone
}

func (ds *Domain) Level(page *domain.Page, term []string) domain.SignalLevel {

	tldURL, err := tld.Parse(page.URL)

	if err != nil {
		return domain.SignalLevelNone
	}

	baseLevel := domain.SignalLevelNone

	// full domain matching
	baseLevel += ds.domainMatch(tldURL, term)

	// hostname matching
	baseLevel += ds.hostnameMatch(tldURL, term)

	return baseLevel
}

package signal

import (
	"crawlquery/node/domain"
	"net/url"
)

// embed:

type Domain struct{}

func (ds *Domain) fullDomainMatch(url *url.URL, term []string) domain.SignalLevel {
	// full domain matching
	for _, t := range term {
		if url.Host == t {
			return domain.SignalLevelMax
		}
	}

	return domain.SignalLevelNone
}

func (ds *Domain) hostnameMatch(url *url.URL, term []string) domain.SignalLevel {
	// hostname matching
	for _, t := range term {
		if url.Hostname() == t {
			return domain.SignalLevelVeryHigh
		}
	}

	return domain.SignalLevelNone
}

func (ds *Domain) Level(page *domain.Page, term []string) domain.SignalLevel {

	baseLevel := domain.SignalLevelNone

	parsedUrl, err := url.Parse(page.URL)
	if err != nil {
		return baseLevel
	}

	// full domain matching
	baseLevel += ds.fullDomainMatch(parsedUrl, term)

	// hostname matching
	baseLevel += ds.hostnameMatch(parsedUrl, term)

	return baseLevel
}

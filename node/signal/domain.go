package signal

import (
	"crawlquery/node/domain"

	tld "github.com/jpillora/go-tld"
)

// embed:

type Domain struct{}

func (Domain) Name() string {
	return "domain"
}

func (ds *Domain) domainMatch(url *tld.URL, term []string) domain.SignalLevel {
	if len(term) == 1 {
		if term[0] == url.Domain {
			if url.Subdomain == "" && url.Path == "" {
				return domain.SignalLevelMax * 1000
			}
		}
	}

	// domain matching
	for _, t := range term {
		if url.Domain == t {
			return domain.SignalLevelVeryHigh
		}
	}

	return domain.SignalLevelNone
}

func (ds *Domain) hostnameMatch(url *tld.URL, term []string) domain.SignalLevel {
	if len(term) == 1 {
		if term[0] == url.Hostname() {
			if url.Subdomain == "" && url.Path == "" {
				return domain.SignalLevelMax * 1000
			}
		}
	}
	// hostname matching
	for _, t := range term {
		if url.Host == t {
			return domain.SignalLevelMax
		}
	}

	return domain.SignalLevelNone
}

func (ds *Domain) Level(page *domain.Page, term []string) (domain.SignalLevel, domain.SignalBreakdown) {

	tldURL, err := tld.Parse(page.URL)

	if err != nil {
		return domain.SignalLevelNone, domain.SignalBreakdown{}
	}

	baseLevel := domain.SignalLevelNone

	// full domain matching
	domainLevel := ds.domainMatch(tldURL, term)

	// hostname matching
	hostname := ds.hostnameMatch(tldURL, term)

	baseLevel += domainLevel
	baseLevel += hostname

	return baseLevel, domain.SignalBreakdown{
		"domain":   domainLevel,
		"hostname": hostname,
	}

}

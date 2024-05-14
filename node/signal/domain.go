package signal

import (
	"crawlquery/node/domain"
	sharedDomain "crawlquery/pkg/domain"
	"net/url"
	"strings"
)

type DomainSignal struct{}

func (ds *DomainSignal) fuzzySearch(host string, terms []string) domain.SignalLevel {
	for _, term := range terms {
		if strings.Contains(strings.ToLower(host), strings.ToLower(term)) {
			// Return different levels based on the term or other conditions
			if len(term) > 5 {
				return domain.SignalLevelHigh
			}
			return domain.SignalLevelMedium
		}
	}
	return domain.SignalLevelNone
}

func (ds *DomainSignal) Level(page *sharedDomain.Page, term []string) domain.SignalLevel {

	parsedUrl, err := url.Parse(page.URL)

	if err != nil {
		return domain.SignalLevelNone
	}

	host := parsedUrl.Hostname()

	// full domain matching
	for _, t := range term {
		if host == t {
			return domain.SignalLevelVeryStrong
		}
	}

	// fuzzy search
	return ds.fuzzySearch(host, term)
}

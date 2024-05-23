package signal

import (
	"crawlquery/node/domain"
	"net/url"
	"strings"
)

type Domain struct{}

func (ds *Domain) fuzzySearch(host string, terms []string) domain.SignalLevel {
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

func (ds *Domain) Level(page *domain.Page, term []string) domain.SignalLevel {

	parsedUrl, err := url.Parse(page.URL)

	if err != nil {
		return domain.SignalLevelNone
	}

	var fullDomainMatch bool
	host := parsedUrl.Host

	// full domain matching
	for _, t := range term {
		if host == t {
			fullDomainMatch = true
			break
		}
	}

	if parsedUrl.Path == "" && fullDomainMatch {
		return domain.SignalLevelMax
	}

	var subdomainMatch bool
	subdomain := strings.Split(host, ".")[0]

	// subdomain matching
	for _, t := range term {
		if subdomain == t {
			subdomainMatch = true
		}
	}

	if subdomainMatch && parsedUrl.Path == "" {
		return domain.SignalLevelVeryHigh
	}

	// fuzzy search
	return ds.fuzzySearch(host, term)
}

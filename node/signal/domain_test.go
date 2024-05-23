package signal

import (
	"crawlquery/node/domain"
	"testing"

	tld "github.com/jpillora/go-tld"
)

func TestDomainSignalDomainMatch(t *testing.T) {
	t.Run("adds a max signal level for a full domain match", func(t *testing.T) {
		tldURL, _ := tld.Parse("http://example.com")
		cases := []struct {
			name  string
			url   *tld.URL
			terms []string
			want  domain.SignalLevel
		}{
			{
				name:  "single term match",
				url:   tldURL,
				terms: []string{"example"},
				want:  domain.SignalLevelVeryHigh,
			},
			{
				name:  "multiple term match",
				url:   tldURL,
				terms: []string{"example", "test"},
				want:  domain.SignalLevelVeryHigh,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {

				ds := &Domain{}
				level := ds.domainMatch(tc.url, tc.terms)

				if level != tc.want {
					t.Errorf("Expected %s, got %f", tc.want, level)
				}
			})
		}
	})
}

func TestDomainSignalHostnameMatch(t *testing.T) {
	t.Run("adds a very high signal level for a hostname match", func(t *testing.T) {
		tldURL, _ := tld.Parse("http://example.com")
		cases := []struct {
			name  string
			url   *tld.URL
			terms []string
			want  domain.SignalLevel
		}{
			{
				name:  "single term match",
				url:   tldURL,
				terms: []string{"example.com"},
				want:  domain.SignalLevelMax,
			},
			{
				name:  "multiple term match",
				url:   tldURL,
				terms: []string{"example.com", "example.com"},
				want:  domain.SignalLevelMax,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {

				ds := &Domain{}
				level := ds.hostnameMatch(tc.url, tc.terms)

				if level != tc.want {
					t.Errorf("Expected %s, got %f", tc.want, level)
				}
			})
		}
	})
}

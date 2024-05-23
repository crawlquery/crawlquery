package signal

import (
	"crawlquery/node/domain"
	"net/url"
	"testing"
)

func TestDomainFullDomainMatch(t *testing.T) {
	t.Run("adds a max signal level for a full domain match", func(t *testing.T) {
		cases := []struct {
			name  string
			url   *url.URL
			terms []string
			want  domain.SignalLevel
		}{
			{
				name:  "single term match",
				url:   &url.URL{Host: "example.com"},
				terms: []string{"example.com"},
				want:  domain.SignalLevelMax,
			},
			{
				name:  "multiple term match",
				url:   &url.URL{Host: "example.com"},
				terms: []string{"example.com", "test"},
				want:  domain.SignalLevelMax,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {

				ds := &Domain{}
				level := ds.fullDomainMatch(tc.url, tc.terms)

				if level != tc.want {
					t.Errorf("Expected %s, got %f", tc.want, level)
				}
			})
		}
	})
}

func TestDomainHostnameMatch(t *testing.T) {
	t.Run("adds a very high signal level for a hostname match", func(t *testing.T) {
		cases := []struct {
			name  string
			url   *url.URL
			terms []string
			want  domain.SignalLevel
		}{
			{
				name:  "single term match",
				url:   &url.URL{Host: "example.com"},
				terms: []string{"example"},
				want:  domain.SignalLevelVeryHigh,
			},
			{
				name:  "multiple term match",
				url:   &url.URL{Host: "google.com"},
				terms: []string{"google", "test"},
				want:  domain.SignalLevelVeryHigh,
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

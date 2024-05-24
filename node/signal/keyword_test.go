package signal

import (
	"crawlquery/node/domain"
	"testing"
)

func TestKeyword(t *testing.T) {
	t.Run("adds a low level signal for a keyword match", func(t *testing.T) {
		cases := []struct {
			name  string
			page  *domain.Page
			terms []string
			want  domain.SignalLevel
		}{
			{
				name: "single term match",
				page: &domain.Page{
					Keywords: []string{"example"},
				},
				terms: []string{"example"},
				want:  domain.SignalLevelMax,
			},
			{
				name: "multiple term match",
				page: &domain.Page{
					Keywords: []string{"example", "test"},
				},
				terms: []string{"example", "test"},
				want:  domain.SignalLevelMax * 2,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {

				ks := &Keyword{}

				level, _ := ks.Level(tc.page, tc.terms)

				if level != tc.want {
					t.Errorf("Expected %s, got %v", tc.want, level)
				}
			})
		}

	})
}

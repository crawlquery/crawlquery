package signal_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/signal"
	"testing"
)

func TestKeywordSignal(t *testing.T) {
	t.Run("adds score for each matching keyword", func(t *testing.T) {
		p := &signal.Keyword{}

		page := &domain.Page{
			Keywords: [][]string{{"ride", "a", "bike"}, {"ride"}, {"bike"}},
		}

		terms := []string{"how", "to", "ride", "a", "bike"}

		level, _ := p.Level(page, terms)

		if level != domain.SignalLevelMedium*3 {
			t.Errorf("Expected %f, got %f", domain.SignalLevelMedium, level)
		}
	})
}

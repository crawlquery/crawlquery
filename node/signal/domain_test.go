package signal_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/signal"
	sharedDomain "crawlquery/pkg/domain"
	"testing"
)

func TestDomain(t *testing.T) {
	t.Run("fuzzySearch", func(t *testing.T) {
		// Arrange
		ds := &signal.DomainSignal{}
		terms := []string{"example"}

		// Act
		level := ds.Level(&sharedDomain.Page{
			URL: "http://example.com",
		}, terms)

		// Assert
		if level != domain.SignalLevelHigh {
			t.Errorf("Expected high level, got %v", level)
		}
	})

	t.Run("Apply", func(t *testing.T) {
		// Arrange
		ds := &signal.DomainSignal{}
		page := &sharedDomain.Page{
			URL: "http://example.com",
		}
		term := []string{"example.com"}

		// Act
		level := ds.Level(page, term)

		// Assert
		if level != domain.SignalLevelVeryStrong {
			t.Errorf("Expected very strong level, got %v", level)
		}
	})
}

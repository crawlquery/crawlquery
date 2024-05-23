package signal_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/signal"
	"testing"
)

func TestDomain(t *testing.T) {
	t.Run("fuzzySearch", func(t *testing.T) {
		// Arrange
		ds := &signal.Domain{}
		terms := []string{"example"}

		// Act
		level := ds.Level(&domain.Page{
			URL: "http://example.com",
		}, terms)

		// Assert
		if level != domain.SignalLevelVeryHigh {
			t.Errorf("Expected very high level, got %v", level)
		}
	})

	t.Run("Level", func(t *testing.T) {
		// Arrange
		ds := &signal.Domain{}
		page := &domain.Page{
			URL: "http://example.com",
		}
		term := []string{"example.com"}

		// Act
		level := ds.Level(page, term)

		// Assert
		if level != domain.SignalLevelMax {
			t.Errorf("Expected max level, got %v", level)
		}
	})

	t.Run("youtube.com example", func(t *testing.T) {
		// Arrange
		ds := &signal.Domain{}
		page := &domain.Page{
			URL: "https://youtube.com",
		}
		term := []string{"youtube.com"}

		// Act
		level := ds.Level(page, term)

		// Assert
		if level != domain.SignalLevelMax {
			t.Errorf("Expected max level, got %v", level)
		}

	})
}

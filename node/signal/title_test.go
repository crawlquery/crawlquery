package signal

import (
	"crawlquery/node/domain"
	"testing"
)

func TestTitleAnyMatch(t *testing.T) {
	t.Run("adds medium signal for each matching term in the title", func(t *testing.T) {

		cases := []struct {
			name   string
			title  string
			terms  []string
			result domain.SignalLevel
		}{
			{
				name:   "two matches",
				title:  "Gmail - Free Storage and Email from Google",
				terms:  []string{"gmail", "google"},
				result: domain.SignalLevelVeryHigh,
			},
			{
				name:   "one match",
				title:  "Gmail - Free Storage and Email from Google",
				terms:  []string{"gmail"},
				result: domain.SignalLevelMedium,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				// Arrange
				ts := &Title{}

				// Act
				level := ts.anyMatch(c.title, c.terms)

				// Assert
				if level != c.result {
					t.Errorf("Expected %f, got %f", c.result, level)
				}
			})
		}

	})
}

func TestTitleFullMatch(t *testing.T) {
	t.Run("returns max level if the term exactly matches the title", func(t *testing.T) {
		cases := []struct {
			name   string
			title  string
			terms  []string
			result domain.SignalLevel
		}{
			{
				name:   "exact match",
				title:  "Gmail",
				terms:  []string{"gmail"},
				result: domain.SignalLevelMax,
			},
			{
				name:   "not an exact match",
				title:  "Gmail - Free Storage and Email from Google",
				terms:  []string{"gmail"},
				result: domain.SignalLevelNone,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				// Arrange
				ts := &Title{}

				// Act
				level := ts.fullMatch(c.title, c.terms)

				// Assert
				if level != c.result {
					t.Errorf("Expected %f, got %f", c.result, level)
				}
			})
		}

	})
}

func TestTitle(t *testing.T) {
	t.Run("cleans title", func(t *testing.T) {
		// Arrange
		ts := &Title{}
		terms := []string{"Google", "Drive"}

		// Act
		level, _ := ts.Level(&domain.Page{
			Title: "Google Drive: Sign-in",
		}, terms)

		// Assert
		if level < domain.SignalLevelMedium {
			t.Errorf("Expected non-zero level, got %f", level)
		}
	})

	t.Run("cleans terms", func(t *testing.T) {
		// Arrange
		ts := &Title{}
		terms := []string{"Google", "Drive:"}

		// Act
		level, _ := ts.Level(&domain.Page{
			Title: "Google Drive: Sign-in",
		}, terms)

		// Assert
		if level < domain.SignalLevelMedium {
			t.Errorf("Expected non-zero level, got %f", level)
		}
	})

	t.Run("applies anyMatch", func(t *testing.T) {
		// Arrange
		ts := &Title{}
		terms := []string{"gmail", "google"}

		// Act
		level, _ := ts.Level(&domain.Page{
			Title: "Gmail - Free Storage and Email from Google",
		}, terms)

		// Assert
		if level != domain.SignalLevelVeryHigh {
			t.Errorf("Expected high level, got %f", level)
		}
	})

	t.Run("applies fullMatch", func(t *testing.T) {
		// Arrange
		ts := &Title{}
		terms := []string{"gmail", "google"}

		// Act
		level, _ := ts.Level(&domain.Page{
			Title: "Gmail Google",
		}, terms)

		// Assert
		if level < domain.SignalLevelHigh {
			t.Errorf("Expected at least high, got %f", level)
		}
	})
}

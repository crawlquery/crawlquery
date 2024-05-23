package domain

type Signal interface {
	Level(page *Page, term []string) SignalLevel
}

type SignalLevel float64

const (
	SignalLevelNone       SignalLevel = 0
	SignalLevelVeryLow    SignalLevel = 1
	SignalLevelLow        SignalLevel = 3
	SignalLevelModerate   SignalLevel = 10
	SignalLevelMedium     SignalLevel = 20
	SignalLevelHigh       SignalLevel = 30
	SignalLevelVeryHigh   SignalLevel = 40
	SignalLevelStrong     SignalLevel = 90
	SignalLevelVeryStrong SignalLevel = 150
	SignalLevelMax        SignalLevel = 1000
)

func (s SignalLevel) String() string {
	switch s {
	case SignalLevelNone:
		return "None"
	case SignalLevelVeryLow:
		return "Very Low"
	case SignalLevelLow:
		return "Low"
	case SignalLevelModerate:
		return "Moderate"
	case SignalLevelMedium:
		return "Medium"
	case SignalLevelHigh:
		return "High"
	case SignalLevelVeryHigh:
		return "Very High"
	case SignalLevelStrong:
		return "Strong"
	case SignalLevelVeryStrong:
		return "Very Strong"
	case SignalLevelMax:
		return "Max"
	default:
		return "Unknown"
	}
}

package domain

type SignalLevel float64

const (
	SignalLevelNone       SignalLevel = 0.0
	SignalLevelVeryLow    SignalLevel = 0.1
	SignalLevelLow        SignalLevel = 0.2
	SignalLevelModerate   SignalLevel = 0.3
	SignalLevelMedium     SignalLevel = 0.4
	SignalLevelHigh       SignalLevel = 0.5
	SignalLevelVeryHigh   SignalLevel = 0.6
	SignalLevelStrong     SignalLevel = 0.7
	SignalLevelVeryStrong SignalLevel = 0.8
	SignalLevelMax        SignalLevel = 1.0
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

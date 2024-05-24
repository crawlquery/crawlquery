package domain

type Result struct {
	PageID  string                     `json:"id"`
	Score   float64                    `json:"score"`
	Page    *ResultPage                `json:"page"`
	Signals map[string]SignalBreakdown `json:"signals"`
}

// Page represents a web page with metadata. Note this does not include the keywords.
type ResultPage struct {
	ID          string `json:"id"`
	Hash        string `json:"hash"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"meta_description"`
}

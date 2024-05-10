package domain

type Result struct {
	ID          string  `json:"id"`
	Url         string  `json:"url"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

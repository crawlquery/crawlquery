package domain

type PageUpdatedEvent struct {
	Page               *Page                         `json:"page"`
	KeywordOccurrences map[Keyword]KeywordOccurrence `json:"keyword_occurrences"`
}

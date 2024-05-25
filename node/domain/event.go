package domain

type PageUpdatedEvent struct {
	Page *Page `json:"page"`
}

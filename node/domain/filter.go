package domain

type FilterService interface {
	PostingsContainsExplicitWords(postings map[string]*Posting) bool
}

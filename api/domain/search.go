package domain

import "crawlquery/pkg/domain"

type SearchService interface {
	Search(term string) ([]domain.Result, error)
}

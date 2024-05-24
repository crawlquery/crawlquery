package domain

type Parser interface {
	Parse(domain *Page) error
}

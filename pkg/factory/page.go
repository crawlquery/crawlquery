package factory

import (
	"crawlquery/pkg/domain"
)

var HomePage = &domain.Page{
	ID:              "1",
	URL:             "https://example.com",
	Title:           "Home",
	MetaDescription: "Welcome to our official website where we offer the latest updates and information.",
}

package factory

import (
	"crawlquery/pkg/domain"
)

func ExampleResults() []domain.Result {
	return []domain.Result{
		{
			ID:          "1b4e28ba-2fa1-11d2-883f-0016d3cca427",
			Url:         "https://www.google.com",
			Title:       "Google",
			Description: "The world's most popular search engine",
			Score:       0.95,
		},
		{
			ID:          "2fa1a8c8-3a0d-45b7-9743-379d9db72bbf",
			Url:         "https://www.bing.com",
			Title:       "Bing",
			Description: "Microsoft's search engine",
			Score:       0.85,
		},
		{
			ID:          "31c2a26f-7b4e-46ef-a64a-123456789012",
			Url:         "https://duckduckgo.com",
			Title:       "DuckDuckGo",
			Description: "Privacy-focused search engine",
			Score:       0.8,
		},
		{
			ID:          "4d76045a-8e4b-474f-82e2-24d45bbba3b8",
			Url:         "https://www.yahoo.com",
			Title:       "Yahoo",
			Description: "Popular search engine and web portal",
			Score:       0.75,
		},
		{
			ID:          "5484b1e2-ec15-4218-9236-07a1f91e9fc5",
			Url:         "https://www.ask.com",
			Title:       "Ask",
			Description: "Q&A focused search engine",
			Score:       0.7,
		},
	}
}

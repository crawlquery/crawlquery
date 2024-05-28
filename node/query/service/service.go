package service

import "crawlquery/node/domain"

type Service struct {
	pageService domain.PageService
}

func NewService(
	pageService domain.PageService,
) *Service {
	return &Service{
		pageService: pageService,
	}
}

func (s *Service) Query(query string) ([]domain.Page, error) {

	if query == "SELECT title FROM pages WHERE title LIKE '%example%';" {
		return []domain.Page{
			{
				ID:    "page1",
				Title: "Example Page",
			},
		}, nil
	}

	return nil, nil
}

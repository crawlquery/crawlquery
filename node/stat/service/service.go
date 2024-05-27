package service

import (
	"crawlquery/node/domain"
)

type Service struct {
	pageService    domain.PageService
	keywordService domain.KeywordService
	dumpService    domain.DumpService
}

func NewService(
	pageService domain.PageService,
	keywordService domain.KeywordService,
	dumpService domain.DumpService,
) *Service {
	return &Service{
		pageService:    pageService,
		keywordService: keywordService,
		dumpService:    dumpService,
	}
}

func (s *Service) Info() (*domain.StatInfo, error) {
	pages, err := s.pageService.GetAll()
	if err != nil {
		return nil, err
	}

	keywordCount, err := s.keywordService.Count()

	if err != nil {
		return nil, err
	}

	totalPages := len(pages)
	totalIndexedPages := 0
	totalKeywords := keywordCount
	sizeOfPages := 0
	bytes, err := s.dumpService.Page()

	if err != nil {
		return nil, err
	}

	sizeOfPages = len(bytes)

	for _, page := range pages {
		if page.LastIndexedAt != nil {
			totalIndexedPages++
		}
	}

	return &domain.StatInfo{
		TotalPages:        totalPages,
		TotalIndexedPages: totalIndexedPages,
		TotalKeywords:     totalKeywords,
		SizeOfPages:       sizeOfPages,
	}, nil
}

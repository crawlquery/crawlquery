package service

import "crawlquery/node/domain"

type Service struct {
	pageService domain.PageService
	dumpService domain.DumpService
}

func NewService(pageService domain.PageService, dumpService domain.DumpService) *Service {
	return &Service{
		pageService: pageService,
		dumpService: dumpService,
	}
}

func (s *Service) Info() (*domain.StatInfo, error) {
	pages, err := s.pageService.GetAll()
	if err != nil {
		return nil, err
	}

	totalPages := len(pages)
	totalPhrases := 0
	sizeOfIndex := 0
	bytes, err := s.dumpService.Page()

	if err != nil {
		return nil, err
	}

	sizeOfIndex = len(bytes)

	for _, page := range pages {
		totalPhrases += len(page.Phrases)
	}

	return &domain.StatInfo{
		TotalPages:   totalPages,
		TotalPhrases: totalPhrases,
		SizeOfIndex:  sizeOfIndex,
	}, nil
}

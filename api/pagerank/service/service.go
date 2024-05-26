package service

import (
	"crawlquery/api/domain"
	"math"
	"time"

	"go.uber.org/zap"
)

const (
	d       = 0.85  // Damping factor
	epsilon = 0.001 // Convergence threshold
	maxIter = 100   // Maximum iterations
)

type Service struct {
	linkService  domain.LinkService
	pageRankRepo domain.PageRankRepository
	logger       *zap.SugaredLogger
}

func NewService(linkService domain.LinkService, pageRankRepo domain.PageRankRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		linkService:  linkService,
		pageRankRepo: pageRankRepo,
		logger:       logger,
	}
}

func (s *Service) UpdatePageRanks() error {
	pages, links, err := s.fetchPagesAndLinks()
	if err != nil {
		return err
	}

	pageRanks := calculatePageRank(pages, links)

	for pageID, rank := range pageRanks {

		err := s.pageRankRepo.Update(pageID, rank, time.Now())
		if err != nil {
			s.logger.Errorw("Error updating page rank", "error", err)
			return err
		}
	}

	return nil
}

func (s *Service) GetPageRank(pageID string) (float64, error) {

	rank, err := s.pageRankRepo.Get(pageID)
	if err != nil {
		s.logger.Errorw("Error getting page rank", "error", err)
		return 0, err
	}

	return rank, nil
}

func (s *Service) fetchPagesAndLinks() (map[string]*domain.PageRank, map[string][]string, error) {
	pages := make(map[string]*domain.PageRank)
	links := make(map[string][]string)

	allLinks, err := s.linkService.GetAll()
	if err != nil {
		s.logger.Errorw("Error getting all links", "error", err)
		return nil, nil, err
	}

	for _, link := range allLinks {
		if _, ok := pages[link.SrcID]; !ok {
			pages[link.SrcID] = &domain.PageRank{
				PageID:   link.SrcID,
				PageRank: 1.0,
			}
		}

		if _, ok := pages[link.DstID]; !ok {
			pages[link.DstID] = &domain.PageRank{
				PageID:   link.DstID,
				PageRank: 1.0,
			}
		}

		links[link.SrcID] = append(links[link.SrcID], link.DstID)
	}

	return pages, links, nil
}

func calculatePageRank(pages map[string]*domain.PageRank, links map[string][]string) map[string]float64 {
	numPages := len(pages)
	pageRanks := make(map[string]float64)
	newPageRanks := make(map[string]float64)

	for id := range pages {
		pageRanks[id] = 1.0 / float64(numPages)
	}

	for i := 0; i < maxIter; i++ {
		for id := range pages {
			newPageRanks[id] = (1.0 - d) / float64(numPages)
		}

		for src, dsts := range links {
			srcRank := pageRanks[src]
			lenDsts := float64(len(dsts))

			for _, dst := range dsts {
				newPageRanks[dst] += d * (srcRank / lenDsts)
			}
		}

		var diff float64
		for id := range pages {
			diff += math.Abs(newPageRanks[id] - pageRanks[id])
			pageRanks[id] = newPageRanks[id]
		}

		if diff < epsilon {
			break
		}
	}

	return pageRanks
}

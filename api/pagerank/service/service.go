package service

import (
	"crawlquery/api/domain"
	nodeDomain "crawlquery/node/domain"
	"crawlquery/pkg/testutil"
	"math"

	"go.uber.org/zap"
)

const (
	d       = 0.85  // Damping factor
	epsilon = 0.001 // Convergence threshold
	maxIter = 100   // Maximum iterations
)

type Service struct {
	linkService domain.LinkService
	logger      *zap.SugaredLogger
}

func NewService(linkService domain.LinkService, logger *zap.SugaredLogger) *Service {
	return &Service{
		linkService: linkService,
		logger:      logger,
	}
}

func (s *Service) ApplyPageRankToResults(results []nodeDomain.Result) ([]nodeDomain.Result, error) {
	pages, links, err := s.fetchPagesAndLinks()
	if err != nil {
		return nil, err
	}

	pageRanks := calculatePageRank(pages, links)

	testutil.PrettyPrint(pageRanks)
	var updatedResults []nodeDomain.Result

	for _, result := range results {
		if rank, ok := pageRanks[result.PageID]; ok {
			result.PageRank = rank
		}
		updatedResults = append(updatedResults, result)
	}

	return updatedResults, nil
}

func (s *Service) CalculatePageRank(pageID string) (float64, error) {

	pages, links, err := s.fetchPagesAndLinks()
	if err != nil {
		s.logger.Errorw("Error fetching pages and links", "error", err)
		return 0, err
	}

	pageRanks := calculatePageRank(pages, links)

	if rank, ok := pageRanks[pageID]; ok {
		return rank, nil
	}

	s.logger.Errorw("Page not found", "pageID", pageID)

	return 0, nil
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

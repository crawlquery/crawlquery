package service

import (
	"crawlquery/node/domain"
	"sort"
	"strings"
)

type Service struct {
	pageService    domain.PageService
	keywordService domain.KeywordService
}

func NewService(
	pageService domain.PageService,
	keywordService domain.KeywordService,
) *Service {
	return &Service{
		pageService:    pageService,
		keywordService: keywordService,
	}
}

func sortResultsByHits(results []*domain.Result) {
	sort.Slice(results, func(i, j int) bool {
		return len(results[i].Hits) > len(results[j].Hits)
	})
}

func (s *Service) getResultsForTermGroup(termGroups [][]string) ([]*domain.Result, error) {
	allPages := make(map[string]*domain.Page)
	hits := make(map[string]map[string]int) // pageID -> keyword -> count

	for _, group := range termGroups {
		groupQuery := strings.Join(group, " ")
		pageIDs, err := s.keywordService.GetPageIDsByKeyword(groupQuery)
		if err != nil {
			continue
		}

		for _, pageID := range pageIDs {
			if _, ok := allPages[pageID]; !ok {
				page, err := s.pageService.Get(pageID)
				if err != nil {
					return nil, err
				}
				allPages[pageID] = page
				hits[pageID] = make(map[string]int)
			}
			hits[pageID][groupQuery]++
		}
	}

	results := []*domain.Result{}
	for pageID, page := range allPages {
		resultPage := &domain.ResultPage{
			ID:          page.ID,
			Hash:        page.Hash,
			URL:         page.URL,
			Title:       page.Title,
			Description: page.Description,
		}
		results = append(results, &domain.Result{
			PageID: pageID,
			Page:   resultPage,
			Hits:   hits[pageID],
			Score:  0, // You can implement scoring logic here if needed
		})
	}

	sortResultsByHits(results)

	return results, nil
}

func (s *Service) Search(query string) ([]*domain.Result, error) {

	queryGroups := splitQueryIntoCombinations(query)

	results, err := s.getResultsForTermGroup(queryGroups)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func splitQueryIntoCombinations(query string) [][]string {
	clean := strings.ToLower(strings.Join(strings.Fields(query), " "))
	if clean == "" {
		return [][]string{}
	}

	terms := strings.Split(clean, " ")
	var groups [][]string

	for i := 0; i < len(terms); i++ {
		for j := i; j < len(terms); j++ {
			groups = append(groups, terms[i:j+1])
		}
	}
	return groups
}

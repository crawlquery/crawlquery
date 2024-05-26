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

func sortResults(results []*domain.Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
}

func (s *Service) getResultsForKeywords(keywords []domain.Keyword) ([]*domain.Result, error) {
	unsortedResults := map[string]*domain.Result{}

	matches, err := s.keywordService.GetKeywordMatches(keywords)
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		for _, occurrence := range match.Occurrences {
			page, err := s.pageService.Get(occurrence.PageID)
			if err != nil {
				return nil, err
			}

			if _, ok := unsortedResults[page.ID]; !ok {
				unsortedResults[page.ID] = &domain.Result{
					Page: &domain.ResultPage{
						ID:    page.ID,
						Hash:  page.Hash,
						URL:   page.URL,
						Title: page.Title,
					},
					Score:             0,
					KeywordOccurences: map[string]domain.KeywordOccurrence{},
				}
			}

			unsortedResults[page.ID].KeywordOccurences[string(match.Keyword)] = occurrence

			unsortedResults[page.ID].Score += float64(occurrence.Frequency)
		}
	}

	results := []*domain.Result{}

	for _, result := range results {
		results = append(results, result)
	}

	sortResults(results)

	return results, nil
}

func (s *Service) Search(query string) ([]*domain.Result, error) {
	queryGroups := splitQueryIntoCombinations(query)

	results, err := s.getResultsForKeywords(queryGroups)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func splitQueryIntoCombinations(query string) []domain.Keyword {
	clean := strings.ToLower(strings.Join(strings.Fields(query), " "))
	if clean == "" {
		return []domain.Keyword{}
	}

	terms := strings.Split(clean, " ")

	keywords := []domain.Keyword{}

	for i := 0; i < len(terms); i++ {
		for j := i + 1; j <= len(terms); j++ {
			keywords = append(keywords, domain.Keyword(strings.Join(terms[i:j], " ")))
		}
	}

	return keywords
}

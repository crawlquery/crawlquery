package quality

import "crawlquery/api/domain"

type LanguageQualityScorer struct{}

func (lqs *LanguageQualityScorer) Score(page *domain.Page) int {
	return 0
}

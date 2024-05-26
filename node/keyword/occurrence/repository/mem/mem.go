package mem

import "crawlquery/node/domain"

type Repository struct {
	occurrences map[domain.Keyword][]domain.KeywordOccurrence
}

func NewRepository() *Repository {
	return &Repository{
		occurrences: make(map[domain.Keyword][]domain.KeywordOccurrence),
	}
}

func (r *Repository) GetAll(keyword domain.Keyword) ([]domain.KeywordOccurrence, error) {
	occurrences, ok := r.occurrences[keyword]
	if !ok {
		return nil, domain.ErrKeywordNotFound
	}

	return occurrences, nil
}

func (r *Repository) Add(keyword domain.Keyword, occurrence domain.KeywordOccurrence) error {
	r.occurrences[keyword] = append(r.occurrences[keyword], occurrence)
	return nil
}

func (r *Repository) RemoveForPageID(pageID string) error {
	for keyword, occurrences := range r.occurrences {
		var newOccurrences []domain.KeywordOccurrence

		for _, occurrence := range occurrences {
			if occurrence.PageID != pageID {
				newOccurrences = append(newOccurrences, occurrence)
			}
		}

		r.occurrences[keyword] = newOccurrences

		if len(r.occurrences[keyword]) == 0 {
			delete(r.occurrences, keyword)
		}
	}

	return nil
}

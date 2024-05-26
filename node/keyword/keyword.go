package keyword

import "crawlquery/node/domain"

func MakeKeywordOccurrences(keywords []domain.Keyword, pageID string) (map[domain.Keyword]domain.KeywordOccurrence, error) {

	// Update occurrences
	keywordOccurrences := make(map[domain.Keyword]domain.KeywordOccurrence, 0)

	for i, keyword := range keywords {

		// Get the current occurrence or initialize a new one
		occurrence, ok := keywordOccurrences[keyword]
		if !ok {
			occurrence = domain.KeywordOccurrence{
				PageID:    pageID,
				Frequency: 1,
				Positions: []int{i},
			}
		} else {
			// Update the existing occurrence
			occurrence.Frequency += 1
			occurrence.Positions = append(occurrence.Positions, i)
		}

		// Put the updated occurrence back into the map
		keywordOccurrences[keyword] = occurrence
	}

	return keywordOccurrences, nil
}

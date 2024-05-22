package service

import (
	"crawlquery/node/domain"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

type Service struct {
	repo domain.KeywordRepository
}

func NewService(repo domain.KeywordRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SavePostings(keywordPostings map[string]*domain.Posting) error {
	for keyword, posting := range keywordPostings {
		if err := s.repo.SavePosting(keyword, posting); err != nil {
			return err
		}

		if err := s.UpdateKeywordHash(keyword); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetPostings(keyword string) ([]*domain.Posting, error) {
	postings, err := s.repo.GetPostings(keyword)
	if err != nil {
		return nil, err
	}

	return postings, nil
}

func (s *Service) FuzzySearch(token string) ([]string, error) {
	results := s.repo.FuzzySearch(token)
	return results, nil
}

func (s *Service) RemovePostingsByPageID(pageID string) error {
	return s.repo.RemovePostingsByPageID(pageID)
}

func (s *Service) SerializePostings(postings []*domain.Posting) ([]byte, error) {
	// Ensure the postings are serialized in a consistent order
	sort.Slice(postings, func(i, j int) bool {
		if postings[i].PageID != postings[j].PageID {
			return postings[i].PageID < postings[j].PageID
		}
		if postings[i].Frequency != postings[j].Frequency {
			return postings[i].Frequency < postings[j].Frequency
		}
		// Compare positions lexicographically
		for k := 0; k < len(postings[i].Positions) && k < len(postings[j].Positions); k++ {
			if postings[i].Positions[k] != postings[j].Positions[k] {
				return postings[i].Positions[k] < postings[j].Positions[k]
			}
		}
		return len(postings[i].Positions) < len(postings[j].Positions)
	})

	// Serialize the sorted postings to JSON
	return json.Marshal(postings)
}

func computeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *Service) UpdateKeywordHash(hash string) error {
	postings, err := s.GetPostings(hash)

	if err != nil {
		return err
	}

	serializedPostings, err := s.SerializePostings(postings)

	if err != nil {
		return err
	}

	newHash := computeHash(serializedPostings)

	if newHash != hash {
		return s.repo.UpdateHash(hash, newHash)
	}

	return nil
}

func (s *Service) Hash() (string, error) {
	// get all keyword hashes
	hashes, err := s.repo.GetHashes()

	if err != nil {
		return "", err
	}

	// Sort the hashes by keyword
	keys := make([]string, 0, len(hashes))
	for key := range hashes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Concatenate the sorted hashes
	var concatenatedHashes string
	for _, key := range keys {
		concatenatedHashes += hashes[key]
	}

	// Compute the global hash from the concatenated hashes
	globalHash := sha256.Sum256([]byte(concatenatedHashes))
	return hex.EncodeToString(globalHash[:]), nil
}

func (s *Service) JSON() ([]byte, error) {
	keywords, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	serializedKeywords, err := json.Marshal(keywords)

	if err != nil {
		return nil, err
	}

	return serializedKeywords, nil
}

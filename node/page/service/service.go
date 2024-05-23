package service

import (
	"crawlquery/node/domain"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

type Service struct {
	pageRepo domain.PageRepository
}

func NewService(pr domain.PageRepository) *Service {
	return &Service{
		pageRepo: pr,
	}
}

func computeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *Service) UpdatePageHash(page *domain.Page) error {
	encoded, err := json.Marshal(page)

	if err != nil {
		return err
	}

	hash := computeHash(encoded)

	return s.pageRepo.UpdateHash(page.ID, hash)
}

func (s *Service) Create(pageID, url, hash string) (*domain.Page, error) {

	page := &domain.Page{
		ID:   pageID,
		URL:  url,
		Hash: hash,
	}

	err := s.pageRepo.Save(pageID, page)

	if err != nil {
		return nil, err
	}

	err = s.UpdatePageHash(page)

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *Service) GetAll() (map[string]*domain.Page, error) {
	return s.pageRepo.GetAll()
}

func (s *Service) Delete(pageID string) error {
	err := s.pageRepo.Delete(pageID)

	if err != nil {
		return err
	}

	return s.pageRepo.DeleteHash(pageID)
}

func (s *Service) Update(page *domain.Page) error {
	return s.pageRepo.Save(page.ID, page)
}

func (s *Service) Get(pageID string) (*domain.Page, error) {
	return s.pageRepo.Get(pageID)
}

func (s *Service) Hash() (string, error) {
	// get all keyword hashes
	hashes, err := s.pageRepo.GetHashes()

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
	pages, err := s.pageRepo.GetAll()

	if err != nil {
		return nil, err
	}

	return json.Marshal(pages)
}

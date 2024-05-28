package service

import (
	"crawlquery/node/domain"
	"crawlquery/pkg/util"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

type Service struct {
	pageRepo    domain.PageRepository
	peerService domain.PeerService
}

func NewService(pr domain.PageRepository, peerService domain.PeerService) *Service {
	return &Service{
		pageRepo:    pr,
		peerService: peerService,
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

func (s *Service) Count() (int, error) {
	return s.pageRepo.Count()
}

func (s *Service) GetByIDs(pageIDs []string) (map[string]*domain.Page, error) {
	return s.pageRepo.GetByIDs(pageIDs)
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

	if s.peerService != nil {
		s.peerService.BroadcastPageUpdatedEvent(&domain.PageUpdatedEvent{
			Page: page,
		})
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
	err := s.pageRepo.Save(page.ID, page)

	if err != nil {
		return err
	}

	if s.peerService != nil {
		s.peerService.BroadcastPageUpdatedEvent(&domain.PageUpdatedEvent{
			Page: page,
		})
	}

	return s.UpdatePageHash(page)
}

func (s *Service) UpdateQuietly(page *domain.Page) error {
	err := s.pageRepo.Save(page.ID, page)

	if err != nil {
		return err
	}

	return s.UpdatePageHash(page)
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

	return util.Sha256Hex32([]byte(concatenatedHashes)), nil
}

func (s *Service) JSON() ([]byte, error) {
	pages, err := s.pageRepo.GetAll()

	if err != nil {
		return nil, err
	}

	return json.Marshal(pages)
}

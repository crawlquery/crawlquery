package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/keyword"
	"crawlquery/node/parse"
	"crawlquery/pkg/util"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type Service struct {
	pageService    domain.PageService
	htmlService    domain.HTMLService
	peerService    domain.PeerService
	keywordService domain.KeywordService
	logger         *zap.SugaredLogger
}

func NewService(
	pageService domain.PageService,
	htmlService domain.HTMLService,
	peerService domain.PeerService,
	keywordService domain.KeywordService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		pageService:    pageService,
		htmlService:    htmlService,
		peerService:    peerService,
		keywordService: keywordService,
		logger:         logger,
	}
}

func (s *Service) Index(pageID string) error {
	page, err := s.pageService.Get(pageID)
	if err != nil {
		s.logger.Errorw("Error getting page", "error", err, "pageID", page, "pageID", pageID)
		return err
	}

	html, err := s.htmlService.Get(page.ID)

	if err != nil {
		s.logger.Errorw("Error getting html", "error", err, "pageID", pageID)
		return err
	}

	page.Hash = util.Sha256Hex32(html)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	if err != nil {
		s.logger.Errorw("Error parsing html", "error", err, "pageID", pageID)
		return err
	}

	title, err := parse.Title(doc)
	if err != nil {
		s.logger.Errorw("Error parsing title", "error", err, "pageID", pageID)
	}

	desc, err := parse.Description(doc)

	if err != nil {
		s.logger.Errorw("Error parsing description", "error", err, "pageID", pageID)
	}

	keywords, err := parse.Keywords(doc)

	if err != nil {
		s.logger.Errorw("Error parsing keywords", "error", err, "pageID", pageID)
	}

	if len(keywords) > 1500 {
		keywords = keywords[:1500]

		s.logger.Warnw("Truncating keywords", "pageID", pageID)
	}

	occurrences, err := keyword.MakeKeywordOccurrences(keywords, page.ID)

	if err != nil {
		s.logger.Errorw("Error making keyword occurrences", "error", err, "pageID", pageID)
	}

	// Update keywords
	err = s.keywordService.UpdateOccurrences(page.ID, occurrences)

	if err != nil {
		s.logger.Errorw("Error updating keyword occurrences", "error", err, "pageID", pageID)
		return err
	}

	page.Title = title
	page.Description = desc
	now := time.Now()
	page.LastIndexedAt = &now

	err = s.pageService.Update(page)

	if err != nil {
		s.logger.Errorw("Error updating page", "error", err, "pageID", pageID)
		return err
	}

	go s.peerService.BroadcastPageUpdatedEvent(&domain.PageUpdatedEvent{
		Page: page,
	})

	return nil
}

func (s *Service) GetIndex(pageID string) (*domain.Page, error) {
	page, err := s.pageService.Get(pageID)

	if err != nil {
		s.logger.Errorw("Error getting page", "error", err, "pageID", page)
		return nil, err
	}

	return page, nil
}

func (s *Service) ApplyPageUpdatedEvent(event *domain.PageUpdatedEvent) error {
	// update the page
	err := s.pageService.UpdateQuietly(event.Page)

	if err != nil {
		s.logger.Errorf("Error updating page: %v", err)
		return err
	}

	return nil
}

func (s *Service) Hash() (string, error) {
	pageHash, err := s.pageService.Hash()

	if err != nil {
		return "", err
	}

	return pageHash, nil
}

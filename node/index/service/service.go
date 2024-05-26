package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"crawlquery/pkg/util"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type Service struct {
	pageService              domain.PageService
	htmlService              domain.HTMLService
	peerService              domain.PeerService
	keywordOccurrenceService domain.KeywordOccurrenceService
	logger                   *zap.SugaredLogger
}

func NewService(
	pageService domain.PageService,
	htmlService domain.HTMLService,
	peerService domain.PeerService,
	keywordOccurrenceService domain.KeywordOccurrenceService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		pageService:              pageService,
		htmlService:              htmlService,
		peerService:              peerService,
		keywordOccurrenceService: keywordOccurrenceService,
		logger:                   logger,
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

	var keywords [][]string

	parsers := []domain.Parser{
		parse.NewLanguageParser(doc),
		parse.NewTitleParser(doc),
		parse.NewDescriptionParser(doc),
		parse.NewKeywordParser(doc, &keywords),
	}

	var errors []error

	for _, parser := range parsers {
		err = parser.Parse(page)

		if err != nil {
			s.logger.Errorw("Error parsing page", "error", err, "pageID", pageID)
			errors = append(errors, err)
		}
	}
	now := time.Now()
	page.LastIndexedAt = &now

	err = s.pageService.Update(page)

	if err != nil {
		s.logger.Errorw("Error saving page", "error", err, "pageID", pageID)
		return err
	}

	if len(errors) > 0 {
		return errors[0]
	}

	// Update occurrences
	keywordOccurrences := make(map[domain.Keyword]domain.KeywordOccurrence, 0)

	for i, termSplit := range keywords {
		term := strings.Join(termSplit, " ")

		// Get the current occurrence or initialize a new one
		occurrence, ok := keywordOccurrences[domain.Keyword(term)]
		if !ok {
			occurrence = domain.KeywordOccurrence{
				PageID:    page.ID,
				Frequency: 1,
				Positions: []int{i},
			}
		} else {
			// Update the existing occurrence
			occurrence.Frequency += 1
			occurrence.Positions = append(occurrence.Positions, i)
		}

		// Put the updated occurrence back into the map
		keywordOccurrences[domain.Keyword(term)] = occurrence
	}

	// Update keywords
	err = s.keywordOccurrenceService.Update(page.ID, keywordOccurrences)

	if err != nil {
		s.logger.Errorw("Error updating keyword occurrences", "error", err, "pageID", pageID)
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

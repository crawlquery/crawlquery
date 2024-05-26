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

	// Update keywords
	err = s.keywordService.UpdatePageKeywords(page.ID, keywords)

	if err != nil {
		s.logger.Errorw("Error updating keywords", "error", err, "pageID", pageID)
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

func (s *Service) keywordPages(terms [][]string) (map[string]*domain.Page, error) {
	var pages map[string]*domain.Page = make(map[string]*domain.Page)

	for _, termSplit := range terms {

		term := strings.Join(termSplit, " ")
		pageIDs, err := s.keywordService.GetPageIDsByKeyword(term)

		if err != nil {
			continue
		}

		for _, pageID := range pageIDs {

			if _, ok := pages[pageID]; ok {
				continue
			}
			page, err := s.pageService.Get(pageID)

			if err != nil {
				s.logger.Errorw("Error getting page", "error", err, "pageID", pageID)
				return nil, err
			}

			pages[pageID] = page
		}
	}

	return pages, nil
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

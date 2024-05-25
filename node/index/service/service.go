package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"crawlquery/node/signal"
	"crawlquery/pkg/util"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type Service struct {
	pageService domain.PageService
	htmlService domain.HTMLService
	peerService domain.PeerService
	logger      *zap.SugaredLogger
}

func NewService(
	pageService domain.PageService,
	htmlService domain.HTMLService,
	peerService domain.PeerService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		pageService: pageService,
		htmlService: htmlService,
		peerService: peerService,
		logger:      logger,
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

	parsers := []domain.Parser{
		parse.NewLanguageParser(doc),
		parse.NewTitleParser(doc),
		parse.NewDescriptionParser(doc),
		parse.NewPhraseParser(doc),
	}

	for _, parser := range parsers {
		err = parser.Parse(page)

		if err != nil {
			s.logger.Errorw("Error parsing page", "error", err, "pageID", pageID)
			return err
		}
	}

	err = s.pageService.Update(page)

	if err != nil {
		s.logger.Errorw("Error saving page", "error", err, "pageID", pageID)
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

func (s *Service) Search(query string) ([]domain.Result, error) {

	term := strings.Split(query, " ")
	// Get all pages
	// Apply signal level to each page
	// Sort by signal level

	pages, err := s.pageService.GetAll()

	if err != nil {
		s.logger.Errorw("Error getting pages", "error", err)
		return nil, err
	}

	results := make([]domain.Result, 0)

	signals := []domain.Signal{
		&signal.Domain{},
		&signal.Title{},
		&signal.Phrase{},
	}

	for _, page := range pages {

		var breakdown map[string]domain.SignalBreakdown = make(map[string]domain.SignalBreakdown)
		totalScore := 0.0
		for _, signal := range signals {
			val, sigs := signal.Level(page, term)
			breakdown[signal.Name()] = sigs
			totalScore += float64(val)
		}

		if totalScore > 1 {
			results = append(results, domain.Result{
				PageID: page.ID,
				Page: &domain.ResultPage{
					ID:          page.ID,
					Hash:        page.Hash,
					URL:         page.URL,
					Title:       page.Title,
					Description: page.Description,
				},
				Signals: breakdown,
				Score:   totalScore,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > 10 {
		results = results[:10]
	}

	return results, nil
}

func (s *Service) ApplyPageUpdatedEvent(event *domain.PageUpdatedEvent) error {
	// update the page
	err := s.pageService.Update(event.Page)

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

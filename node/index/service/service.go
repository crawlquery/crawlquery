package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"crawlquery/node/signal"
	"crawlquery/node/token"
	"sort"
	"strings"

	sharedDomain "crawlquery/pkg/domain"
	"crawlquery/pkg/util"

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

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	if err != nil {
		s.logger.Errorw("Error parsing html", "error", err, "pageID", pageID)
		return err
	}

	page.Title = parse.Title(doc)
	page.MetaDescription = parse.MetaDescription(doc)
	page.Hash = util.Sha256Hex32(html)

	// FOR NOW KEYWORDS MUST COME LAST AS IT REMOVES HTML TAGS
	// TO GET THE KEYWORDS
	page.Keywords = token.Keywords(doc)

	if page.MetaDescription == "" {
		subKeyWords := page.Keywords

		if len(subKeyWords) > 10 {
			subKeyWords = subKeyWords[:10]
		}

		page.MetaDescription = strings.Join(subKeyWords, ", ")
	}

	err = s.pageService.Update(page)

	if err != nil {
		s.logger.Errorw("Error saving page", "error", err, "pageID", pageID)
		return err
	}

	go s.peerService.BroadcastIndexEvent(&domain.IndexEvent{
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

func (s *Service) ReIndex(pageID string) error {
	page, err := s.pageService.Get(pageID)

	if err != nil {
		s.logger.Errorw("Error getting page", "error", err, "pageID", page)
		return err
	}

	err = s.Index(pageID)

	if err != nil {
		s.logger.Errorw("Error indexing page", "error", err, "pageID", pageID)
		return err
	}

	return nil
}

func (s *Service) Search(query string) ([]sharedDomain.Result, error) {

	term := strings.Split(query, " ")
	// Get all pages
	// Apply signal level to each page
	// Sort by signal level

	pages, err := s.pageService.GetAll()

	if err != nil {
		s.logger.Errorw("Error getting pages", "error", err)
		return nil, err
	}

	results := make([]sharedDomain.Result, 0)

	signals := []domain.Signal{
		&signal.Domain{},
		&signal.Title{},
		&signal.Keyword{},
	}

	for _, page := range pages {
		score := 0.0
		for _, signal := range signals {
			score += float64(signal.Level(page, term))
		}

		if score > 1 {
			results = append(results, sharedDomain.Result{
				PageID: page.ID,
				Page: &sharedDomain.Page{
					ID:              page.ID,
					Hash:            page.Hash,
					URL:             page.URL,
					Title:           page.Title,
					MetaDescription: page.MetaDescription,
				},
				Score: score,
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

func (s *Service) ApplyIndexEvent(event *domain.IndexEvent) error {
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

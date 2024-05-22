package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"crawlquery/node/signal"
	"crawlquery/node/token"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	sharedDomain "crawlquery/pkg/domain"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type Service struct {
	pageService    domain.PageService
	htmlService    domain.HTMLService
	keywordService domain.KeywordService
	peerService    domain.PeerService
	logger         *zap.SugaredLogger
}

func NewService(
	pageService domain.PageService,
	htmlService domain.HTMLService,
	keywordService domain.KeywordService,
	peerService domain.PeerService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		pageService:    pageService,
		htmlService:    htmlService,
		keywordService: keywordService,
		peerService:    peerService,
		logger:         logger,
	}
}

func (s *Service) MakePostings(page *sharedDomain.Page, keywords []string) map[string]*domain.Posting {
	postings := make(map[string]*domain.Posting, 0)

	for i, keyword := range keywords {
		if _, ok := postings[keyword]; !ok {
			postings[keyword] = &domain.Posting{
				Frequency: 1,
				PageID:    page.ID,
				Positions: []int{i},
			}
		} else {
			postings[keyword].Frequency++
			postings[keyword].Positions = append(postings[keyword].Positions, i)
		}
	}

	return postings
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

	// FOR NOW KEYWORDS MUST COME LAST AS IT REMOVES HTML TAGS
	// TO GET THE KEYWORDS
	keywords := token.Keywords(doc)
	postings := s.MakePostings(page, keywords)

	err = s.keywordService.SavePostings(postings)

	if err != nil {
		s.logger.Errorw("Error getting keywords", "error", err, "pageID", pageID)
		return err
	}

	err = s.pageService.Update(page)

	if err != nil {
		s.logger.Errorw("Error saving page", "error", err, "pageID", pageID)
		return err
	}

	go s.peerService.BroadcastIndexEvent(&domain.IndexEvent{
		Page:     page,
		Keywords: postings,
	})

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
	}

	for _, page := range pages {
		score := 0.0
		for _, signal := range signals {
			score += float64(signal.Level(page, term))
		}

		if score > 1 {
			results = append(results, sharedDomain.Result{
				PageID: page.ID,
				Page:   page,
				Score:  score,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func (s *Service) ApplyIndexEvent(event *domain.IndexEvent) error {
	page, _ := s.pageService.Get(event.Page.ID)
	var err error
	if page == nil {
		page, err = s.pageService.Create(event.Page.ID, event.Page.URL)

		if err != nil {
			s.logger.Errorf("Error creating page: %v", err)
			return err
		}
	}

	// add the title and meta description
	err = s.pageService.Update(event.Page)

	if err != nil {
		s.logger.Errorf("Error updating page: %v", err)
		return err
	}

	// remove old postings
	err = s.keywordService.RemovePostingsByPageID(page.ID)

	if err != nil {
		s.logger.Errorf("Error removing postings: %v", err)
		return err
	}

	// add new postings
	err = s.keywordService.SavePostings(event.Keywords)

	if err != nil {
		s.logger.Errorf("Error saving postings: %v", err)
		return err
	}

	return nil
}

func computeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *Service) Hash() (string, string, string, error) {
	pageHash, err := s.pageService.Hash()

	if err != nil {
		return "", "", "", err
	}

	keywordHash, err := s.keywordService.Hash()

	if err != nil {
		return "", "", "", err
	}

	return pageHash, keywordHash, computeHash([]byte(pageHash + keywordHash)), nil
}

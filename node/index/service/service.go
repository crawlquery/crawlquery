package service

import (
	"bytes"
	"crawlquery/node/domain"
	"crawlquery/node/parse"
	"crawlquery/node/token"
	"crypto/sha256"
	"encoding/hex"
	"sort"

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
	// Tokenize the query the same way as the index was tokenized
	queryTokens := token.TokenizeTerm(query)
	results := make(map[string]float64) // map[PageID]relevanceScore

	// use full token search first
	for _, term := range queryTokens {
		postings, err := s.keywordService.GetPostings(term)

		if err != nil {
			s.logger.Errorf("Error getting postings: %v", err)
			continue
		}

		for _, posting := range postings {
			results[posting.PageID] += float64(posting.Frequency)
		}

	}

	if len(results) == 0 {
		for _, term := range queryTokens {
			fuzzyTokens, err := s.keywordService.FuzzySearch(term)
			if err != nil {
				s.logger.Errorf("Error getting fuzzy search results: %v", err)
				continue
			}

			for _, token := range fuzzyTokens {
				postings, err := s.keywordService.GetPostings(token)
				if err != nil {
					s.logger.Errorf("Error getting postings: %v", err)
					continue
				}
				for _, posting := range postings {
					results[posting.PageID] += float64(posting.Frequency)
				}
			}
		}
	}

	// Convert the results map to a slice and sort by relevance score
	var sortedResults []sharedDomain.Result
	for docID, score := range results {
		sortedResults = append(sortedResults, sharedDomain.Result{PageID: docID, Score: score})
	}

	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].Score > sortedResults[j].Score
	})

	// Add the page metadata to the results
	for i, result := range sortedResults {
		page, err := s.pageService.Get(result.PageID)

		if err != nil {
			s.logger.Errorf("Index.Search: Error getting page metadata: %v", err)
			continue
		}
		sortedResults[i].Page = page
	}

	if len(sortedResults) >= 10 {
		sortedResults = sortedResults[:10]
	}

	return sortedResults, nil
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

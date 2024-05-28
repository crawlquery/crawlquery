package service

import (
	"crawlquery/node/domain"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repairJobRepo  domain.RepairJobRepository
	pageService    domain.PageService
	keywordService domain.KeywordService
	peerService    domain.PeerService
	logger         *zap.SugaredLogger
}

func NewService(
	repairJobRepo domain.RepairJobRepository,
	pageService domain.PageService,
	keywordService domain.KeywordService,
	peersService domain.PeerService,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		repairJobRepo:  repairJobRepo,
		pageService:    pageService,
		keywordService: keywordService,
		peerService:    peersService,
		logger:         logger,
	}
}

func (s *Service) AuditAndRepair() error {
	currentPages, err := s.pageService.GetAll()

	if err != nil {
		s.logger.Errorw("Error getting all index metas", "error", err)
		return err
	}

	peerMetas, err := s.peerService.GetAllIndexMetas()

	if err != nil {
		s.logger.Errorw("Error getting peer index metas", "error", err)
		return err
	}

	latestIndexedAtPeers := s.MapLatestPages(peerMetas, currentPages)

	peerPages := s.GroupPageIDsByThePeerID(latestIndexedAtPeers)

	var pageIDs []string

	for _, ids := range peerPages {
		for _, id := range ids {
			pageIDs = append(pageIDs, string(id))
		}
	}

	err = s.ProcessRepairJobs(pageIDs)

	if err != nil {
		s.logger.Errorw("Error processing repair jobs", "error", err)
		return err
	}

	return nil
}

func (s *Service) AuditAndRepairEvery(interval int) {
	err := s.AuditAndRepair()

	if err != nil {
		s.logger.Errorw("Error auditing and repairing", "error", err)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := s.AuditAndRepair()

		if err != nil {
			s.logger.Errorw("Error auditing and repairing", "error", err)
		}
	}
}

func (s *Service) GetAllIndexMetas() ([]domain.IndexMeta, error) {
	var metas []domain.IndexMeta

	pages, err := s.pageService.GetAll()

	if err != nil {
		s.logger.Errorw("Error getting pages", "error", err)
		return nil, err
	}

	for _, page := range pages {
		if page.LastIndexedAt == nil {
			continue
		}
		metas = append(metas, domain.IndexMeta{
			PageID:        domain.PageID(page.ID),
			PeerID:        domain.PeerID(s.peerService.Self().ID),
			LastIndexedAt: *page.LastIndexedAt,
		})
	}

	return metas, nil
}

func (s *Service) GetIndexMetas(pageIDs []string) ([]domain.IndexMeta, error) {
	var metas []domain.IndexMeta

	pages, err := s.pageService.GetByIDs(pageIDs)

	if err != nil {
		s.logger.Errorw("Error getting pages", "error", err)
		return nil, err
	}

	for _, page := range pages {
		if page.LastIndexedAt == nil {
			continue
		}
		metas = append(metas, domain.IndexMeta{
			PageID:        domain.PageID(page.ID),
			PeerID:        domain.PeerID(s.peerService.Self().ID),
			LastIndexedAt: *page.LastIndexedAt,
		})
	}

	return metas, nil
}

func (s *Service) GetPageDumps(pageIDs []string) ([]domain.PageDump, error) {
	var dumps []domain.PageDump

	for _, pageID := range pageIDs {
		page, err := s.pageService.Get(pageID)

		if err != nil {
			s.logger.Errorw("Error getting page", "error", err)
			return nil, err
		}

		keywordOccurrences, err := s.keywordService.GetForPageID(pageID)

		if err != nil {
			s.logger.Errorw("Error getting keyword occurrences", "error", err)
			return nil, err
		}

		dumps = append(dumps, domain.PageDump{
			PeerID:             domain.PeerID(s.peerService.Self().ID),
			PageID:             domain.PageID(page.ID),
			Page:               *page,
			KeywordOccurrences: keywordOccurrences,
		})
	}

	return dumps, nil
}

func (s *Service) CreateRepairJobs(pageIDs []string) error {
	for _, pageID := range pageIDs {
		_, err := s.repairJobRepo.Get(pageID)

		if err == domain.ErrRepairJobNotFound {
			err = s.repairJobRepo.Create(&domain.RepairJob{
				PageID: pageID,
			})
		}

		if err != nil {
			s.logger.Errorw("Error creating repair job", "error", err)
			return err
		}
	}

	return nil
}

func (s *Service) MapLatestPages(metas []domain.IndexMeta, currentPages map[string]*domain.Page) domain.LatestIndexedPages {
	latestIndexedAtPeers := make(domain.LatestIndexedPages)

	for _, meta := range metas {

		if currentPage, ok := currentPages[string(meta.PageID)]; ok {
			if currentPage.LastIndexedAt != nil && currentPage.LastIndexedAt.After(meta.LastIndexedAt) {
				continue
			}
		}

		if latestIndexedAtPeers[meta.PageID].LatestIndexedAt.Before(meta.LastIndexedAt) {
			latestIndexedAtPeers[meta.PageID] = domain.PeerWithLatestIndexedAt{
				PeerID:          meta.PeerID,
				LatestIndexedAt: meta.LastIndexedAt,
			}
		}
	}

	return latestIndexedAtPeers
}

func (s *Service) GroupPageIDsByThePeerID(latestIndexedAtPeers domain.LatestIndexedPages) domain.PeerPages {
	peerIDToPageIDs := make(domain.PeerPages)

	for pageID, peerWithLatestIndexedAt := range latestIndexedAtPeers {
		peerIDToPageIDs[peerWithLatestIndexedAt.PeerID] = append(peerIDToPageIDs[peerWithLatestIndexedAt.PeerID], pageID)
	}

	return peerIDToPageIDs
}

func (s *Service) ProcessRepairJobs(pageIDs []string) error {

	currentPages, err := s.pageService.GetByIDs(pageIDs)

	if err != nil {
		s.logger.Errorw("Error getting current pages", "error", err)
		return err
	}

	metas, err := s.peerService.GetIndexMetas(pageIDs)

	if err != nil {
		s.logger.Errorw("Error getting index metas", "error", err)
		return err
	}

	latestIndexedAtPeers := s.MapLatestPages(metas, currentPages)

	peerPages := s.GroupPageIDsByThePeerID(latestIndexedAtPeers)

	var dumps []domain.PageDump

	for peerID, pageIDs := range peerPages {
		peer, err := s.peerService.GetPeer(string(peerID))

		if err != nil {
			s.logger.Errorw("Error getting peer", "error", err)
			return err
		}

		peerDumps, err := s.peerService.GetPageDumpsFromPeer(peer, pageIDs)

		if err != nil {
			s.logger.Errorw("Error getting page dumps from peer", "error", err)
			return err
		}

		dumps = append(dumps, peerDumps...)
	}

	for _, dump := range dumps {

		page := dump.Page
		err := s.pageService.UpdateQuietly(&page)

		if err != nil {
			s.logger.Errorw("Error updating page", "error", err)
			return err
		}

		err = s.keywordService.UpdateOccurrences(page.ID, dump.KeywordOccurrences)

		if err != nil {
			s.logger.Errorw("Error updating keyword occurrences", "error", err)
			return err
		}
	}

	return nil
}

package service

import (
	"crawlquery/node/domain"

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

func (s *Service) MapLatestPages(metas []domain.IndexMeta) domain.LatestIndexedPages {
	latestIndexedAtPeers := make(domain.LatestIndexedPages)

	for _, meta := range metas {

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

	metas, err := s.peerService.GetIndexMetas(pageIDs)

	if err != nil {
		s.logger.Errorw("Error getting index metas", "error", err)
		return err
	}

	latestIndexedAtPeers := s.MapLatestPages(metas)

	peerIDToPageIDs := s.GroupPageIDsByThePeerID(latestIndexedAtPeers)

	var dumps []*domain.PageDump

	for peerID, pageIDs := range peerIDToPageIDs {
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
		page := &dump.Page

		err := s.pageService.UpdateQuietly(page)

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

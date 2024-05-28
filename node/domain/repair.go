package domain

import (
	"crawlquery/node/dto"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrRepairJobNotFound = errors.New("repair job not found")

type RepairJobStatus string

const (
	RepairJobStatusPending  RepairJobStatus = "pending"
	RepairJobStatusRunning  RepairJobStatus = "running"
	RepairJobStatusComplete RepairJobStatus = "complete"
	RepairJobStatusFailed   RepairJobStatus = "failed"
)

type RepairJob struct {
	PageID              string
	Status              RepairJobStatus
	StatusLastUpdatedAt time.Time
	CreatedAt           time.Time
}

type RepairService interface {
	CreateRepairJobs(pageID []string) error
	GetAllIndexMetas() ([]IndexMeta, error)
	GetIndexMetas(pageIDs []string) ([]IndexMeta, error)
	GetPageDumps(pageIDs []string) ([]PageDump, error)
}

type RepairHandler interface {
	GetIndexMetas(c *gin.Context)
	GetAllIndexMetas(c *gin.Context)
	GetPageDumps(c *gin.Context)
}

type RepairJobRepository interface {
	Create(pageID *RepairJob) error
	Get(pageID string) (*RepairJob, error)
	Update(job *RepairJob) error
}

type PeerWithLatestIndexedAt struct {
	PeerID          PeerID
	LatestIndexedAt time.Time
}

type LatestIndexedPages map[PageID]PeerWithLatestIndexedAt

type PeerPages map[PeerID][]PageID

type IndexMeta struct {
	PageID        PageID
	PeerID        PeerID
	LastIndexedAt time.Time
}

type PageDump struct {
	PeerID             PeerID
	PageID             PageID
	Page               Page
	KeywordOccurrences map[Keyword]KeywordOccurrence
}

func PageDumpFromDTO(d dto.PageDump) PageDump {
	pageDump := PageDump{
		PeerID: PeerID(d.PeerID),
		PageID: PageID(d.PageID),
		Page: Page{
			ID:            d.Page.ID,
			URL:           d.Page.URL,
			Title:         d.Page.Title,
			Description:   d.Page.Description,
			Language:      d.Page.Language,
			Hash:          d.Page.Hash,
			LastIndexedAt: &d.Page.LastIndexedAt,
		},
		KeywordOccurrences: make(map[Keyword]KeywordOccurrence),
	}

	for keyword, occurrence := range d.KeywordOccurrences {
		pageDump.KeywordOccurrences[Keyword(keyword)] = KeywordOccurrence{
			PageID:    occurrence.PageID,
			Frequency: occurrence.Frequency,
			Positions: occurrence.Positions,
		}
	}

	return pageDump
}

func PageDumpToDTO(d PageDump) dto.PageDump {

	lastIndexedAt := time.Time{}

	if d.Page.LastIndexedAt != nil {
		lastIndexedAt = *d.Page.LastIndexedAt
	}
	newDto := dto.PageDump{
		PeerID: string(d.PeerID),
		PageID: string(d.PageID),
		Page: dto.Page{
			ID:            d.Page.ID,
			URL:           d.Page.URL,
			Title:         d.Page.Title,
			Description:   d.Page.Description,
			Language:      d.Page.Language,
			Hash:          d.Page.Hash,
			LastIndexedAt: lastIndexedAt,
		},
		KeywordOccurrences: make(map[string]dto.KeywordOccurrence),
	}

	for keyword, occurrence := range d.KeywordOccurrences {
		newDto.KeywordOccurrences[string(keyword)] = dto.KeywordOccurrence{
			PageID:    occurrence.PageID,
			Frequency: occurrence.Frequency,
			Positions: occurrence.Positions,
		}
	}

	return newDto
}

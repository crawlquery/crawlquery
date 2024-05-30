package mem

import (
	"crawlquery/api/domain"
	"sync"
)

type Repository struct {
	logs  map[domain.CrawlLogID]*domain.CrawlLog
	mutex *sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		logs:  make(map[domain.CrawlLogID]*domain.CrawlLog),
		mutex: &sync.Mutex{},
	}
}

func (r *Repository) Save(cl *domain.CrawlLog) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.logs[cl.ID] = cl
	return nil
}

func (r *Repository) ListByPageID(pageID domain.PageID) ([]*domain.CrawlLog, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var logs []*domain.CrawlLog
	for _, log := range r.logs {
		if log.PageID == pageID {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

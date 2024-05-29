package mem

import "crawlquery/api/domain"

type Repository struct {
	logs map[domain.CrawlLogID]*domain.CrawlLog
}

func NewRepository() *Repository {
	return &Repository{
		logs: make(map[domain.CrawlLogID]*domain.CrawlLog),
	}
}

func (r *Repository) Save(cl *domain.CrawlLog) error {
	r.logs[cl.ID] = cl
	return nil
}

func (r *Repository) ListByPageID(pageID domain.PageID) ([]*domain.CrawlLog, error) {
	var logs []*domain.CrawlLog
	for _, log := range r.logs {
		if log.PageID == pageID {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

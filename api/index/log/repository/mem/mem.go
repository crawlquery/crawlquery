package mem

import (
	"crawlquery/api/domain"
)

type Repository struct {
	logs map[domain.IndexLogID]*domain.IndexLog
}

func NewRepository() *Repository {
	return &Repository{
		logs: make(map[domain.IndexLogID]*domain.IndexLog),
	}
}

func (r *Repository) Save(cl *domain.IndexLog) error {
	r.logs[cl.ID] = cl
	return nil
}

func (r *Repository) ListByPageID(pageID domain.PageID) ([]*domain.IndexLog, error) {
	var logs []*domain.IndexLog
	for _, log := range r.logs {
		if log.PageID == pageID {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

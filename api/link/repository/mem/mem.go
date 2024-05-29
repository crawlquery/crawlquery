package mem

import (
	"crawlquery/api/domain"
)

type Repository struct {
	links []*domain.Link
}

func NewRepository() *Repository {
	return &Repository{
		links: []*domain.Link{},
	}
}

func (r *Repository) Create(link *domain.Link) error {
	for _, l := range r.links {
		if l.SrcID == link.SrcID && l.DstID == link.DstID {
			return domain.ErrLinkAlreadyExists
		}
	}

	r.links = append(r.links, link)

	return nil
}

func (r *Repository) GetAll() ([]*domain.Link, error) {
	return r.links, nil
}

func (r *Repository) GetAllBySrcID(srcID domain.PageID) ([]*domain.Link, error) {
	var links []*domain.Link

	for _, l := range r.links {
		if l.SrcID == srcID {
			links = append(links, l)
		}
	}

	return links, nil
}

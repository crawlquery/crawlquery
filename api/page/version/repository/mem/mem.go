package mem

import "crawlquery/api/domain"

type Repository struct {
	versions map[domain.PageVersionID]*domain.PageVersion
}

func NewRepository() *Repository {
	return &Repository{
		versions: make(map[domain.PageVersionID]*domain.PageVersion),
	}
}

func (r *Repository) Get(id domain.PageVersionID) (*domain.PageVersion, error) {
	version, ok := r.versions[id]
	if !ok {
		return nil, domain.ErrPageVersionNotFound
	}
	return version, nil
}

func (r *Repository) Create(v *domain.PageVersion) error {
	r.versions[v.ID] = v
	return nil
}

func (r *Repository) ListByPageID(pageID domain.PageID) ([]*domain.PageVersion, error) {
	var versions []*domain.PageVersion
	for _, version := range r.versions {
		if version.PageID == pageID {
			versions = append(versions, version)
		}
	}
	return versions, nil
}

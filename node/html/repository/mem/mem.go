package mem

import "crawlquery/node/domain"

type Repository struct {
	htmls map[string][]byte
}

func NewRepository() *Repository {
	return &Repository{
		htmls: make(map[string][]byte),
	}
}

func (r *Repository) Save(pageID string, html []byte) error {
	r.htmls[pageID] = html
	return nil
}

func (r *Repository) Get(pageID string) ([]byte, error) {
	html, ok := r.htmls[pageID]
	if !ok {
		return nil, domain.ErrHTMLNotFound
	}
	return html, nil
}

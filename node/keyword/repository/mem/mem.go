package mem

type Repository struct {
	keywords map[string][]string
}

func NewRepository() *Repository {
	return &Repository{
		keywords: make(map[string][]string),
	}
}

func (r *Repository) GetPages(keyword string) ([]string, error) {
	return r.keywords[keyword], nil
}

func (r *Repository) AddPageKeywords(pageID string, keywords []string) error {
	for _, keyword := range keywords {
		r.keywords[keyword] = append(r.keywords[keyword], pageID)
	}

	return nil
}

func (r *Repository) RemovePageKeywords(pageID string) error {
	for keyword, pages := range r.keywords {
		for i, page := range pages {
			if page == pageID {
				r.keywords[keyword] = append(pages[:i], pages[i+1:]...)
			}
		}
	}

	return nil
}

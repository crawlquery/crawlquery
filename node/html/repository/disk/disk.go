package disk

import (
	"crawlquery/node/domain"
	"fmt"
	"io"
	"os"
)

type Repository struct {
	path string
}

func NewRepository(path string) (*Repository, error) {

	err := os.MkdirAll(path, 0755)

	if err != nil {
		return nil, err
	}

	return &Repository{
		path: path,
	}, nil
}

func (r *Repository) Save(pageID string, data []byte) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", r.path, pageID))

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(pageID string) ([]byte, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", r.path, pageID))

	if err != nil {
		return nil, domain.ErrHTMLNotFound
	}

	defer file.Close()

	return io.ReadAll(file)
}

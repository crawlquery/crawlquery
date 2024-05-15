package bolt

import (
	"crawlquery/pkg/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type Repository struct {
	db *bolt.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	// Ensure the bucket exists
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("ForwardIndex"))
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return &Repository{db: db}, nil
}

func (repo *Repository) Save(pageID string, page *domain.Page) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ForwardIndex"))
		encoded, err := json.Marshal(page)
		if err != nil {
			return err
		}
		return b.Put([]byte(pageID), encoded)
	})
}

func (repo *Repository) Get(pageID string) (*domain.Page, error) {
	var page *domain.Page
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ForwardIndex"))
		v := b.Get([]byte(pageID))
		if v == nil {
			return fmt.Errorf("page not found")
		}
		return json.Unmarshal(v, &page)
	})
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

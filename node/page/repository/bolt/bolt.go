package bolt

import (
	"crawlquery/node/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type Repository struct {
	db *bolt.DB
}

var pageBucket = []byte("Pages")
var hashBucket = []byte("Hashes")

func NewRepository(dbPath string) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	// Ensure the bucket exists
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(pageBucket)

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(hashBucket)

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
		b := tx.Bucket(pageBucket)
		encoded, err := json.Marshal(page)
		if err != nil {
			return err
		}
		return b.Put([]byte(pageID), encoded)
	})
}

func (repo *Repository) Delete(pageID string) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(pageBucket)
		return b.Delete([]byte(pageID))
	})
}

func (repo *Repository) Get(pageID string) (*domain.Page, error) {
	var page *domain.Page
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(pageBucket)
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

func (repo *Repository) Count() (int, error) {
	count := 0
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(pageBucket)
		count = b.Stats().KeyN
		return nil
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *Repository) GetAll() (map[string]*domain.Page, error) {
	pages := make(map[string]*domain.Page)
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(pageBucket)
		return b.ForEach(func(k, v []byte) error {
			page := &domain.Page{}
			err := json.Unmarshal(v, page)
			if err != nil {
				return err
			}
			pages[string(k)] = page
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (repo *Repository) GetByIDs(pageIDs []string) (map[string]*domain.Page, error) {
	pages := make(map[string]*domain.Page)

	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(pageBucket)
		for _, pageID := range pageIDs {
			v := b.Get([]byte(pageID))
			if v == nil {
				return fmt.Errorf("page not found")
			}
			page := &domain.Page{}
			err := json.Unmarshal(v, page)
			if err != nil {
				return err
			}
			pages[pageID] = page
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return pages, nil
}

func (repo *Repository) UpdateHash(pageID string, hash string) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(hashBucket)
		return b.Put([]byte(pageID), []byte(hash))
	})
}

func (repo *Repository) DeleteHash(pageID string) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(hashBucket)
		return b.Delete([]byte(pageID))
	})
}

func (repo *Repository) GetHash(pageID string) (string, error) {
	var hash string
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(hashBucket)
		v := b.Get([]byte(pageID))
		if v == nil {
			return fmt.Errorf("hash not found")
		}
		hash = string(v)
		return nil
	})
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (repo *Repository) GetHashes() (map[string]string, error) {
	hashes := make(map[string]string)
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(hashBucket)
		return b.ForEach(func(k, v []byte) error {
			hashes[string(k)] = string(v)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return hashes, nil
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

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

func NewRepository(dbPath string) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	// Ensure the bucket exists
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("InvertedIndex"))
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return &Repository{db: db}, nil
}

func (repo *Repository) Save(keyword string, postings []*domain.Posting) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("InvertedIndex"))
		encoded, err := json.Marshal(postings)
		if err != nil {
			return err
		}
		return b.Put([]byte(keyword), encoded)
	})
}

func (repo *Repository) Get(keyword string) ([]*domain.Posting, error) {
	var postings []*domain.Posting
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("InvertedIndex"))
		v := b.Get([]byte(keyword))
		if v == nil {
			return nil // Not found
		}
		return json.Unmarshal(v, &postings)
	})
	if err != nil {
		return nil, err
	}
	return postings, nil
}

func (repo *Repository) FuzzySearch(token string) map[string]float64 {
	results := make(map[string]float64)

	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("InvertedIndex"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := string(k)
			postings := make([]*domain.Posting, 0)
			err := json.Unmarshal(v, &postings)
			if err != nil {
				return err
			}

			if len(postings) == 0 {
				continue
			}

			if len(token) > len(key) {
				continue
			}

			if key == token {
				for _, posting := range postings {
					results[posting.PageID] += float64(posting.Frequency)
				}
			} else {
				for i := 0; i < len(key)-len(token)+1; i++ {
					if key[i:i+len(token)] == token {
						for _, posting := range postings {
							results[posting.PageID] += float64(posting.Frequency)
						}
						break
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return results
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

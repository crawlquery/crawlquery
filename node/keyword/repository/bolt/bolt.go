package bolt

import (
	"crawlquery/node/domain"
	"encoding/json"
	"fmt"
	"strings"
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

func (repo *Repository) SavePosting(keyword string, posting *domain.Posting) error {
	err := repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("InvertedIndex"))
		v := b.Get([]byte(keyword))
		if v == nil {
			postings := []*domain.Posting{posting}
			encoded, err := json.Marshal(postings)
			if err != nil {
				return err
			}
			return b.Put([]byte(keyword), encoded)
		}
		var existing []*domain.Posting
		err := json.Unmarshal(v, &existing)
		if err != nil {
			return err
		}
		existing = append(existing, posting)
		encoded, err := json.Marshal(existing)
		if err != nil {
			return err
		}
		return b.Put([]byte(keyword), encoded)
	})
	return err
}

func (repo *Repository) GetPostings(keyword string) ([]*domain.Posting, error) {
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

func (repo *Repository) FuzzySearch(token string) []string {
	results := []string{}

	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("InvertedIndex"))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			key := string(k)
			if strings.Contains(key, token) {
				results = append(results, key)
			}
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return results
}

func (repo *Repository) RemovePostingsByPageID(pageID string) error {
	err := repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("InvertedIndex"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var postings []*domain.Posting
			err := json.Unmarshal(v, &postings)
			if err != nil {
				return err
			}
			for i, posting := range postings {
				if posting.PageID == pageID {
					postings = append(postings[:i], postings[i+1:]...)
				}
			}
			encoded, err := json.Marshal(postings)
			if err != nil {
				return err
			}
			err = b.Put(k, encoded)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

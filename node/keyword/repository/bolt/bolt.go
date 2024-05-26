package bolt

import (
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
		_, err := tx.CreateBucketIfNotExists([]byte("Keywords"))

		return err
	})

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return &Repository{db: db}, nil
}

func (r *Repository) AddPageKeywords(pageID string, keywords []string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keywords"))

		for _, keyword := range keywords {
			encoded, err := json.Marshal([]string{pageID})
			// If the keyword already exists, append the pageID to the list
			if val := b.Get([]byte(keyword)); val != nil {
				var pages []string
				err = json.Unmarshal(val, &pages)
				if err != nil {
					return err
				}
				pages = append(pages, pageID)
				encoded, err = json.Marshal(pages)
				if err != nil {
					return err
				}
				continue
			}
			err = b.Put([]byte(keyword), encoded)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) GetPages(keyword string) ([]string, error) {
	var pages []string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keywords"))
		v := b.Get([]byte(keyword))
		if v == nil {
			return fmt.Errorf("keyword not found")
		}
		return json.Unmarshal(v, &pages)
	})
	return pages, err
}

func (r *Repository) RemovePageKeywords(pageID string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keywords"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var pages []string
			err := json.Unmarshal(v, &pages)
			if err != nil {
				return err
			}
			for i, page := range pages {
				if page == pageID {
					pages = append(pages[:i], pages[i+1:]...)
					encoded, err := json.Marshal(pages)
					if err != nil {
						return err
					}
					err = b.Put(k, encoded)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

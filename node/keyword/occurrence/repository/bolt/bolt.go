package bolt

import (
	"crawlquery/node/domain"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

var occurrencesBucket = []byte("occurrences")

type Repository struct {
	db *bolt.DB
}

func NewRepository(db *bolt.DB) (*Repository, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(occurrencesBucket)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetAll(keyword domain.Keyword) ([]domain.KeywordOccurrence, error) {
	var occurrences []domain.KeywordOccurrence

	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(occurrencesBucket)
		if bucket == nil {
			return domain.ErrKeywordNotFound
		}

		data := bucket.Get([]byte(keyword))
		if data == nil {
			return domain.ErrKeywordNotFound
		}

		return json.Unmarshal(data, &occurrences)
	})

	if err != nil {
		return nil, err
	}

	return occurrences, nil
}

func (r *Repository) Add(keyword domain.Keyword, occurrence domain.KeywordOccurrence) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(occurrencesBucket)
		if bucket == nil {
			return errors.New("bucket not found")
		}

		var occurrences []domain.KeywordOccurrence
		data := bucket.Get([]byte(keyword))
		if data != nil {
			err := json.Unmarshal(data, &occurrences)
			if err != nil {
				return err
			}
		}

		occurrences = append(occurrences, occurrence)
		encoded, err := json.Marshal(occurrences)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(keyword), encoded)
	})
}

func (r *Repository) RemoveForPageID(pageID string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(occurrencesBucket)
		if bucket == nil {
			return errors.New("bucket not found")
		}

		err := bucket.ForEach(func(k, v []byte) error {
			var occurrences []domain.KeywordOccurrence
			err := json.Unmarshal(v, &occurrences)
			if err != nil {
				return err
			}

			var newOccurrences []domain.KeywordOccurrence
			for _, occ := range occurrences {
				if occ.PageID != pageID {
					newOccurrences = append(newOccurrences, occ)
				}
			}

			if len(newOccurrences) > 0 {
				encoded, err := json.Marshal(newOccurrences)
				if err != nil {
					return err
				}
				return bucket.Put(k, encoded)
			}

			return bucket.Delete(k)
		})

		return err
	})
}

func (r *Repository) Count() (int, error) {
	var count int

	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(occurrencesBucket)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			count++
			return nil
		})
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}
